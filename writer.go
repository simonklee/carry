// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package carry

import (
	"fmt"
	"sync"
	"time"

	"github.com/simonz05/carry/types"
)

type StatsWriter interface {
	Write(stats []*types.Stat) error
}

type statsBuffer struct {
	w   StatsWriter
	buf []*types.Stat
}

func newStatsBuffer(w StatsWriter) *statsBuffer {
	return &statsBuffer{
		w:   w,
		buf: make([]*types.Stat, 0),
	}
}

func (b *statsBuffer) Flush() error {
	if b.Len() == 0 {
		return nil
	}

	err := b.w.Write(b.buf)

	if err != nil {
		return err
	}

	b.buf = b.buf[0:0]
	return nil
}

func (b *statsBuffer) Len() int {
	return len(b.buf)
}

func (b *statsBuffer) Write(p []*types.Stat) error {
	b.buf = append(b.buf, p...)
	return nil
}

type PeriodicWriter struct {
	buf    *statsBuffer
	in     chan []*types.Stat
	wg     sync.WaitGroup
	period time.Duration
}

func NewPeriodicWriter(w StatsWriter) *PeriodicWriter {
	pw := &PeriodicWriter{
		buf:    newStatsBuffer(w),
		in:     make(chan []*types.Stat, 1024),
		period: (100 * time.Millisecond),
	}

	pw.wg.Add(1)
	go func() {
		pw.process()
		pw.wg.Done()
	}()

	return pw
}

func (pw *PeriodicWriter) process() {
	var lastWrite time.Time

	for {
		now := time.Now()
		dt := pw.period - now.Sub(lastWrite)

		select {
		case stats := <-pw.in:
			pw.buf.Write(stats)
		case <-time.After(dt):
			pw.buf.Flush()
			lastWrite = time.Now()
		}
	}
}

func (pw *PeriodicWriter) Close() error {
	close(pw.in)
	done := make(chan bool)

	func(done chan bool) {
		pw.wg.Wait()
		done <- true
	}(done)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		return fmt.Errorf("Timed out")
	}
	return nil
}

func (pw *PeriodicWriter) Write(stats []*types.Stat) error {
	pw.in <- stats
	return nil
}
