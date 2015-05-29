// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import "github.com/BurntSushi/toml"

type Config struct {
	Listen      string
	Periodic    bool
	AllowOrigin []string `toml:"allow-origin"`
	Stathat     *Stathat
	InfluxDB    *InfluxDB `toml:"influxdb"`
}

type Stathat struct {
	Key string
}

type InfluxDB struct {
	Host     string
	Password string
	Username string
	Database string
}

func ReadFile(filename string) (*Config, error) {
	config := new(Config)
	_, err := toml.DecodeFile(filename, config)

	if len(config.AllowOrigin) == 0 {
		config.AllowOrigin = []string{"*"}
	}

	return config, err
}
