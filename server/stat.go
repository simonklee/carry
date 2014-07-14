// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"net/http"

	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/httputil"
	"github.com/simonz05/util/log"
)

func (c *context) createStat(rw http.ResponseWriter, req *http.Request) {
	if req.ContentLength > MaxStatSize {
		httputil.ServeJSONCodeError(rw, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	stats := []*types.Stat{}
	d := json.NewDecoder(req.Body)
	err := d.Decode(&stats)

	if err != nil {
		log.Error(err)
		httputil.ServeJSONCodeError(rw, "JSON Decode Error", http.StatusBadGateway)
		return
	}

	err = c.sto.ReceiveStats(stats)

	if err != nil {
		log.Error(err)
		httputil.ServeJSONCodeError(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (c *context) headStat(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}
