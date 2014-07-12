// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package storagetest

import (
	"github.com/simonz05/carry"
	"github.com/simonz05/carry/types"
)

type fakeStorage struct {
	stats []*types.Stat
}

func NewFakeStorage() carry.Storage {
	return &fakeStorage{
		stats: make([]*types.Stat, 0),
	}
}

func (sto *fakeStorage) ReceiveStats(stats []*types.Stat) error {
	sto.stats = append(sto.stats, stats...)
	return nil
}
