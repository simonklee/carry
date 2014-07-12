// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

import (
	"fmt"
	"time"
)

type StatKind int

const (
	_                    = iota
	CounterKind StatKind = iota
	ValueKind
)

type Stat struct {
	Key       string   `json:"k"`
	Value     float64  `json:"v"`
	Timestamp int64    `json:"t"`
	Type      StatKind `json:"c"`
}

func (s *Stat) String() string {
	return fmt.Sprintf("stat: %s, status: %f, timestamp: %v, type: %d",
		s.Key, s.Value, time.Unix(s.Timestamp, 0).UTC(), s.Type)
}

func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0 // return correctly abs(-0)
	}
	return x
}
