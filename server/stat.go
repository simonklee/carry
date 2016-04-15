// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
			if s.Timestamp <= 0 {
				s.Timestamp = now
			}

			if s.Value < 0 {
				s.Value = 0
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

func (c *context) createStatGet(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}

	keys, ok := req.Form["k"]
	n := len(keys)

	if !ok || n == 0 {
		httputil.BadRequestError(rw, "No stats")
		return
	}

	stats := make([]*types.Stat, 0, len(keys))

	for _, key := range keys {
		stats = append(stats, &types.Stat{Key: key})
	}

	values, ok := req.Form["v"]

	if !ok || n != len(values) {
		httputil.BadRequestError(rw, fmt.Sprintf("bad value part ok %v, len %d", ok, len(values)))
		return
	}

	for i, value := range values {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			httputil.BadRequestError(rw, fmt.Sprintf("bad value part %s %s", value, err))
			return
		}

		if v < 0 {
			v = 0.0
		}

		stats[i].Value = v
	}

	times, ok := req.Form["t"]

	if !ok || n != len(times) {
		httputil.BadRequestError(rw, "bad timestamp part")
		return
	}

	for i, ts := range times {
		t, err := strconv.ParseInt(ts, 10, 64)
		if err != nil {
			httputil.BadRequestError(rw, "bad timestamp part")
			return
		}

		if t < 0 {
			t = 0
		}

		stats[i].Timestamp = t
	}

	kinds, ok := req.Form["c"]

	if !ok || n != len(kinds) {
		httputil.BadRequestError(rw, "bad stat kind part")
		return
	}

	for i, value := range kinds {
		kind := new(types.StatKind)
		err := kind.UnmarshalText([]byte(value))

		if err != nil {
			httputil.BadRequestError(rw, "bad stat kind part")
			return
		}
		stats[i].Type = *kind
	}

	err = c.sto.ReceiveStats(stats)

	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Internal Server Error"))
		return
	}

	rw.Header().Add("Cache-Control", "private, no-cache, no-cache=Set-Cookie, proxy-revalidate")
	rw.Header().Add("Pragma", "no-cache")
	rw.WriteHeader(http.StatusNoContent)
}

func (c *context) headStat(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Methods", "PUT, DELETE, GET, POST, HEAD, OPTIONS")
	rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, *")
	//rw.Header().Add("Access-Control-Allow-Headers", "origin, x-csrftoken, content-type, accept")
	rw.WriteHeader(http.StatusOK)
}
