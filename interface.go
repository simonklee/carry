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
