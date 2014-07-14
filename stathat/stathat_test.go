// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package stathat

import (
	"os"
	"testing"

	"github.com/simonz05/carry"
	"github.com/simonz05/carry/config"
	"github.com/simonz05/carry/storagetest"
	"github.com/simonz05/util/log"
)

func readConf(t *testing.T) *config.Config {
	configFile := os.Getenv("STATHAT_TEST_CONFIG")
	if configFile == "" {
		t.Skip("Skipping manual test. To enable, set the environment variable STATHAT_TEST_CONFIG to the path of a configuration for the storage type.")
	}
	conf, err := config.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Error reading stathat configuration file %s: %v", configFile, err)
	}
	return conf
}

func TestStathat(t *testing.T) {
	if testing.Verbose() {
		log.Severity = log.LevelInfo
	} else {
		log.Severity = log.LevelError
	}
	conf := readConf(t)
	storagetest.Test(t, func(t *testing.T) (sto carry.Storage, cleanup func()) {
		sto, err := newFromConfig(conf)

		if err != nil {
			t.Fatalf("newFromConfig error: %v", err)
		}

		closer, ok := sto.(carry.ShutdownStorage)

		if !ok {
			t.Fatalf("expected stathat shutdown storage")
		}

		return sto, func() {
			log.Println("cleanup stathat")
			if err := closer.Close(); err != nil {
				t.Fatal(err)
			}
		}
	})
}
