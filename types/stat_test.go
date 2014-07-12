// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	res := &Stat{
		Key:       "k",
		Value:     3.14,
		Timestamp: time.Now().Unix(),
		Type:      ValueKind,
	}

	b, err := json.MarshalIndent(res, "  ", "")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%s - %s", b, res)
}
