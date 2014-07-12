// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package carry

import (
	"fmt"
	"sync"

	"github.com/simonz05/carry/config"
)

// A StorageConstructor returns a Storage implementation from a configuration.
type StorageConstructor func(*config.Config) (Storage, error)

var mapLock sync.Mutex
var storageConstructors = make(map[string]StorageConstructor)

func RegisterStorageConstructor(typ string, ctor StorageConstructor) {
	mapLock.Lock()
	defer mapLock.Unlock()
	if _, ok := storageConstructors[typ]; ok {
		panic("StorageConstructor already registered for type: " + typ)
	}
	storageConstructors[typ] = ctor
}

func CreateStorage(typ string, config *config.Config) (Storage, error) {
	mapLock.Lock()
	ctor, ok := storageConstructors[typ]
	mapLock.Unlock()
	if !ok {
		return nil, fmt.Errorf("Storage type %s not known or loaded", typ)
	}
	return ctor(config)
}
