// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package carry implements a service for processing carrys.

package carry

import (
	"io"

	"github.com/simonz05/carry/types"
)

// StatsReceiver is the interface for receiving
type StatsReceiver interface {
	// ReceiveStats accepts a carry and writes it to
	// storage.
	ReceiveStats([]*types.Stat) error
}

// Storage is the interface that must be implemented by a
// storage type.
type Storage interface {
	StatsReceiver
}

// Optional interface for storage implementations which can be asked
// to shut down cleanly. Regardless, all implementations should
// be able to survive crashes without data loss.
type ShutdownStorage interface {
	Storage
	io.Closer
}

type Loader interface {
	GetStorage() (Storage, error)
}

// MultiStorage implements the storage interface for many storage backends.
type MultiStorage []Storage

// ReceiveStats implements StatsReceiver.
func (ms MultiStorage) ReceiveStats(stats []*types.Stat) (err error) {
	for _, s := range ms {
		if err1 := s.ReceiveStats(stats); err == nil && err1 != nil {
			err = err1
		}
	}
	return
}

// ReceiveStats implements io.Closer.
func (ms MultiStorage) Close() (err error) {
	for _, s := range ms {
		cl, ok := s.(ShutdownStorage)

		if !ok {
			continue
		}

		if err1 := cl.Close(); err == nil && err1 != nil {
			err = err1
		}
	}
	return
}
