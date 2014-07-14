package stathat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/simonz05/carry/types"
	"github.com/simonz05/util/log"
)

type StathatWriter struct {
	key string
}

func NewStathatWriter(key string) *StathatWriter {
	return &StathatWriter{key: key}
}

func (sw *StathatWriter) Write(stats []*types.Stat) error {
	data := make([]ProtocolStat, len(stats))

	for i, stat := range stats {
		data[i].Stat = stat.Key
		data[i].Timestamp = stat.Timestamp
		switch stat.Type {
		case types.ValueKind:
			data[i].Value = &stat.Value
		case types.CounterKind:
			v := int64(stat.Value)
			data[i].Count = &v
		}
	}

	pkg := &ProtocolPackage{
		EZKey: sw.key,
		Data:  data,
	}

	var rw bytes.Buffer
	enc := json.NewEncoder(&rw)
	err := enc.Encode(pkg)

	if err != nil {
		return err
	}

	buflen := rw.Len()
	req, err := http.NewRequest("POST", "http://api.stathat.com/ez", &rw)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(buflen))
	res, err := http.DefaultClient.Do(req)
	log.Println(res.Request.URL)
	return err
}
