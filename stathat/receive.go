// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stathat

import (
	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/log"
)

// ReceiveStats implements the interface for receiving stats.
func (s *stathatStorage) ReceiveStats(stats []*types.Stat) error {
	log.Printf("stathat: receive %d stats", len(stats))
	return s.w.Write(stats)
}
