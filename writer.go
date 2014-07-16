// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package carry

import (
	"fmt"
	"sync"
	"time"

	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/log"
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
		period: (1 * time.Second),
	}

	pw.wg.Add(1)
	go func() {
		pw.process()
		pw.wg.Done()
		log.Println("pw.wg.done")
	}()

	return pw
}

func (pw *PeriodicWriter) process() {
	lastWrite := time.Now()

	for {
		now := time.Now()
		dt := pw.period - now.Sub(lastWrite)
		select {
		case stats, ok := <-pw.in:
			if !ok {
				log.Println("PeriodicWriter chan closed, flush … exit")
				err := pw.buf.Flush()
				if err != nil {
					log.Error(err)
				}
				return
			}

			err := pw.buf.Write(stats)

			if err != nil {
				log.Error(err)
			}
		case <-time.After(dt):
			//log.Printf("PeriodicWriter writing after %v", time.Since(lastWrite))
			err := pw.buf.Flush()
			if err != nil {
				log.Error(err)
			}
			lastWrite = time.Now()
		}
	}
}

func (pw *PeriodicWriter) Close() error {
	log.Println("Closing PeriodicWriter")
	close(pw.in)
	done := make(chan bool)

	go func(done chan bool) {
		pw.wg.Wait()
		done <- true
	}(done)

	select {
	case <-done:
		break
	case <-time.After(1 * time.Second):
		return fmt.Errorf("Timed out")
	}
	log.Println("PeriodicWriter Closed")
	return nil
}

func (pw *PeriodicWriter) Write(stats []*types.Stat) error {
	pw.in <- stats
	return nil
}
