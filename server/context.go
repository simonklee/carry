// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"fmt"

	"github.com/simonz05/carry"
	"github.com/simonz05/carry/config"
	"github.com/simonz05/util/log"
)

type context struct {
	sto         carry.Storage
	allowOrigin []string
}

func newContextFromConfig(conf *config.Config) (*context, error) {
	// We just have stathat storage for now.
	if conf.Stathat == nil && conf.InfluxDB == nil {
		return nil, fmt.Errorf("expected stathat/influxdb storage")
	}

	var storage carry.Storage
	var backends []carry.Storage

	if conf.Stathat != nil {
		sto, err := carry.CreateStorage("stathat", conf)
		if err != nil {
			return nil, err
		}
		backends = append(backends, sto)
		log.Println("stathat storage initialized")
	}

	if conf.InfluxDB != nil {
		sto, err := carry.CreateStorage("influxdb", conf)

		if err != nil {
			return nil, err
		}

		backends = append(backends, sto)
		log.Println("influxdb storage initialized")
	}

	if len(backends) == 1 {
		storage = backends[0]
	} else {
		storage = carry.MultiStorage(backends)
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
