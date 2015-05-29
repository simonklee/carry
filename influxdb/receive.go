// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package influxdb

import (
	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/log"
)

// ReceiveStats implements the interface for receiving stats.
func (s *influxdbStorage) ReceiveStats(stats []*types.Stat) error {
	log.Printf("influxdb: receive %d stats", len(stats))
	return s.w.Write(stats)
}
