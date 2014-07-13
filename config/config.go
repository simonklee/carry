// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Listen   string
	Periodic bool
	Stathat  *Stathat
	Graphite *Graphite
}

type Stathat struct {
	Key string
}

type Graphite struct {
	DSN string `toml:"dsn"`
}

func ReadFile(filename string) (*Config, error) {
	config := new(Config)
	_, err := toml.DecodeFile(filename, config)
	return config, err
}
