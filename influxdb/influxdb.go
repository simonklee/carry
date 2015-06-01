// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package influxdb registers the "influxdb" storage type, storing stats
// in StatHat.

package influxdb

import (
	"fmt"
	"net/url"

	"github.com/influxdb/influxdb/client"

	"github.com/simonz05/carry"
	"github.com/simonz05/carry/config"
)

type influxdbStorage struct {
	config client.Config
	client *client.Client
	w      carry.StatsWriter
}

func (s *influxdbStorage) String() string {
	return fmt.Sprintf("\"influxdb\" storage for %q", s.config.URL)
}

func newFromConfig(conf *config.Config) (carry.Storage, error) {
	host, err := url.Parse(conf.InfluxDB.URL)
	if err != nil {
		return nil, err
	}

	cConfig := client.Config{
		URL:      *host,
		Password: conf.InfluxDB.Password,
		Username: conf.InfluxDB.Username,
	}

	c, err := client.NewClient(cConfig)

	if err != nil {
		return nil, err
	}

	var w carry.StatsWriter
	w = NewInfluxDBWriter(c, conf.InfluxDB.Database)

	if conf.Periodic {
		w = carry.NewPeriodicWriter(w)
	}

	return &influxdbStorage{
		config: cConfig,
		client: c,
		w:      w,
	}, nil
}

func init() {
	carry.RegisterStorageConstructor("influxdb", carry.StorageConstructor(newFromConfig))
}

// compile check to verify influxdb implements carry.ShutdownStorage
var _ carry.ShutdownStorage = &influxdbStorage{}
