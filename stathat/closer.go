// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stathat

import "io"

// Close cleanly shutdown the Stathat
func (s *stathatStorage) Close() error {
	closer, ok := s.w.(io.Closer)
	if ok {
		return closer.Close()
	}
	return nil
}
