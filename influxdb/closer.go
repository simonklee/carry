// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package influxdb

import (
	"io"

	"github.com/simonz05/util/log"
)

// Close cleanly shutdown the InfluxDB
func (s *influxdbStorage) Close() error {
	closer, ok := s.w.(io.Closer)
	if ok {
		log.Println("close influxdb")
		return closer.Close()
	}
	return nil
}
