// Copyright 2015 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package influxdb

import (
	"encoding/json"
	"time"

	"github.com/influxdb/influxdb/client"
	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/log"
)

type InfluxDBWriter struct {
	client *client.Client
	db     string
}

func NewInfluxDBWriter(c *client.Client, db string) *InfluxDBWriter {
	return &InfluxDBWriter{client: c, db: db}
}

func (iw *InfluxDBWriter) Write(stats []*types.Stat) error {
	points := make([]client.Point, len(stats))

	// todo group similar keys.
	for i, stat := range stats {
		points[i] = client.Point{
			Name: stat.Key,
			Fields: map[string]interface{}{
				"value": stat.Value,
			},
			Time: time.Unix(stat.Timestamp, 0),
		}
	}

	if log.Severity >= log.LevelInfo {
		b, _ := json.Marshal(points)
		log.Println("sending stats: ", string(b))
	}

	packet := client.BatchPoints{
		Points:   points,
		Database: iw.db,
	}
	_, err := iw.client.Write(packet)
	return err
}
