package storagetest

import (
	"reflect"
	"testing"
	"time"

	"github.com/simonz05/carry"
	"github.com/simonz05/carry/types"
)

func Test(t *testing.T, fn func(*testing.T) (sto carry.Storage, cleanup func())) {
	sto, cleanup := fn(t)
	defer func() {
		if t.Failed() {
			t.Logf("test %T FAILED, skipping cleanup!", sto)
		} else {
			cleanup()
		}
	}()
	t.Logf("Testing stats storage %T", sto)

	now := time.Now().UTC()
	tests := []struct {
		n string
		s []*types.Stat
		e error
	}{
		{
			n: "single-ok",
			s: []*types.Stat{
				{
					Key:       "k",
					Value:     3.14,
					Timestamp: now.Unix(),
					Type:      types.ValueKind,
				},
			},
		},
		{
			n: "multi-ok",
			s: []*types.Stat{
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

	t.Logf("Testing receive")

	for _, tt := range tests {
		err := sto.ReceiveStats(tt.s)

		if err != nil {
			if tt.e == nil {
				t.Fatalf("[%s]: ReceiveStats of %v", tt.n, err)
			}
			if !reflect.DeepEqual(tt.e, err) {
				t.Fatalf("[%s]: exp err %v got %v", tt.e, err)
			}
			continue
		}
	}
}
