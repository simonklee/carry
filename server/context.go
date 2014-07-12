// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"net/http"

	"github.com/simonz05/carry"
	"github.com/simonz05/carry/config"
	"github.com/simonz05/util/httputil"
)

type context struct {
	sto carry.Storage
}

func newContextFromConfig(conf *config.Config) (*context, error) {
	storage, err := carry.CreateStorage("todo", conf)

	if err != nil {
		return nil, err
	}

	return &context{sto: storage}, nil
}

func (c *context) Close() error {
	closer, ok := c.sto.(carry.ShutdownStorage)
	if ok {
		return closer.Close()
	}
	return nil
}

func (c *context) handlerFunc(f func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := f(w, r)

		if err != nil {
			httputil.ServeJSONError(w, err)
		} else {
			httputil.ReturnJSON(w, data)
		}
	})
}
