// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package server implements a HTTP interface for the social service.

package server

import (
	"io"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/simonz05/carry/config"
	"github.com/simonz05/util/handler"
	"github.com/simonz05/util/ioutil"
	"github.com/simonz05/util/log"
	"github.com/simonz05/util/pat"
	"github.com/simonz05/util/sig"
)

// MaxStatSize is the max size of a stats request
const MaxStatSize = 2 << 20

func Init(conf *config.Config) (closer io.Closer, err error) {
	c, err := newContextFromConfig(conf)
	if err != nil {
		return nil, err
	}
	err = installHandlers(c)
	closer = c
	return
}

func installHandlers(c *context) error {
	router := mux.NewRouter()
	router.StrictSlash(false)

	// package path
	sub := router.PathPrefix("/v1/stat").Subrouter()

	pat.Post(sub, "/p/", http.HandlerFunc(c.createStat))

	// global middleware
	var middleware []func(http.Handler) http.Handler

	if log.Severity >= log.LevelInfo {
		middleware = append(middleware, handler.LogHandler, handler.MeasureHandler, handler.DebugHandle, handler.RecoveryHandler)
	} else {
		middleware = append(middleware, handler.LogHandler, handler.RecoveryHandler)
	}

	wrapped := handler.Use(router, middleware...)
	http.Handle("/", wrapped)
	return nil
}

func ListenAndServe(laddr string, closer io.Closer) error {
	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return err
	}

	log.Printf("Listen on %s", l.Addr())

	closer = ioutil.MultiCloser([]io.Closer{l, closer})
	sig.TrapCloser(closer)
	err = http.Serve(l, nil)
	log.Printf("Shutting down ..")
	return err
}
