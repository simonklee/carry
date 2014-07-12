package server

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"

	"github.com/simonz05/carry/storagetest"
	"github.com/simonz05/carry/types"
	httputil "github.com/simonz05/social/server/httptest"
	"github.com/simonz05/util/assert"
	"github.com/simonz05/util/log"
)

var (
	once       sync.Once
	server     *httptest.Server
	serverAddr string
)

func startServer() {
	ctx := &context{sto: storagetest.NewFakeStorage()}
	err := installHandlers(ctx)

	if err != nil {
		panic(err)
	}

	if testing.Verbose() {
		log.Severity = log.LevelInfo
	} else {
		log.Severity = log.LevelError
	}
	server = httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
}

func TestServer(t *testing.T) {
	once.Do(startServer)
	ast := assert.NewAssertWithName(t, "TestServer")

	tests := []struct {
		n string
		b interface{}
		e error
	}{
		{
			n: "single-ok",
			b: []*types.Stat{
				{
					Key:       "k",
					Value:     3.14,
					Timestamp: time.Now().Unix(),
					Type:      types.ValueKind,
				},
			},
		},
		{
			n: "multi-ok",
			b: []*types.Stat{
				{
					Key:       "k",
					Value:     3.14,
					Timestamp: time.Now().Unix(),
					Type:      types.ValueKind,
				},
				{
					Key:       "k",
					Value:     1.618,
					Timestamp: time.Now().Unix(),
					Type:      types.ValueKind,
				},
			},
		},
	}

	for _, tt := range tests {
		var buf []byte
		if tt.b != nil {
			b, err := json.Marshal(tt.b)
			ast.Nil(err)
			buf = b
		}
		req, err := httputil.NewRequest("POST", absURL("/v1/stat/p/", nil), buf, nil)
		ast.Nil(err)
		req.Header.Set("Content-Type", "application/json")
		res, err := req.Do()

		ast.Nil(err)
		ast.Equal(201, res.StatusCode)
		ast.Nil(err)
	}
}

func absURL(endpoint string, args url.Values) string {
	var params string

	if args != nil && len(args) > 0 {
		params = fmt.Sprintf("?%s", args.Encode())
	}

	return fmt.Sprintf("http://%s%s%s", serverAddr, endpoint, params)
}
