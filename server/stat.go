// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"net/http"
	"time"

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
		httputil.ServeJSONCodeError(rw, "JSON Decode Error", http.StatusBadRequest)
		return
	}

	if len(stats) > 0 {
		now := time.Now().Unix()

		for _, s := range stats {
			if s.Timestamp == 0 {
				s.Timestamp = now
			}
		}
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
	rw.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, POST, HEAD, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, *")
	//rw.Header().Add("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept")
	rw.WriteHeader(http.StatusOK)
}
