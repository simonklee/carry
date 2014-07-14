// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"

	"github.com/simonz05/carry"
	"github.com/simonz05/carry/config"
)

type context struct {
	sto         carry.Storage
	allowOrigin []string
}

func newContextFromConfig(conf *config.Config) (*context, error) {
	// TODO: just have stathat storage for now.
	if conf.Stathat == nil {
		return nil, fmt.Errorf("expected stathat storage")
	}
	storage, err := carry.CreateStorage("stathat", conf)

	if err != nil {
		return nil, err
	}

	return &context{
		sto:         storage,
		allowOrigin: conf.AllowOrigin,
	}, nil
}

func (c *context) Close() error {
	closer, ok := c.sto.(carry.ShutdownStorage)
	if ok {
		return closer.Close()
	}
	return nil
}
