// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package influxdb

import (
	"encoding/json"

	"github.com/influxdb/influxdb/client"
	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/log"
)

type InfluxDBWriter struct {
	client *client.Client
}

func NewInfluxDBWriter(c *client.Client) *InfluxDBWriter {
	return &InfluxDBWriter{client: c}
}

func (iw *InfluxDBWriter) Write(stats []*types.Stat) error {
	data := make([]*client.Series, len(stats))

	// todo group similar keys.
	for i, stat := range stats {
		data[i] = &client.Series{
			Columns: []string{"time", "value"},
			Name:    stat.Key,
			Points: [][]interface{}{
				{stat.Timestamp, stat.Value},
			},
		}
	}

	if log.Severity >= log.LevelInfo {
		b, _ := json.Marshal(data)
		log.Println("sending stats: ", string(b))
	}

	return iw.client.WriteSeries(data)
}
