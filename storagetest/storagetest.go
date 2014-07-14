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
					Key:       "key.value",
					Value:     3.14,
					Timestamp: now.Unix(),
					Type:      types.ValueKind,
				},
				{
					Key:       "key.cnt",
					Value:     2,
					Timestamp: now.Unix(),
					Type:      types.CounterKind,
				},
			},
		},
		{
			n: "multi-ok",
			s: []*types.Stat{
				{
					Key:       "key.value",
					Value:     3.14,
					Timestamp: now.Unix() + 1,
					Type:      types.ValueKind,
				},
				{
					Key:       "key.value",
					Value:     1.618,
					Timestamp: now.Unix() + 2,
					Type:      types.ValueKind,
				},
				{
					Key:       "key.value",
					Value:     8.50,
					Timestamp: now.Unix() + 3,
					Type:      types.ValueKind,
				},
				{
					Key:       "key.cnt",
					Value:     1,
					Timestamp: now.Unix(),
					Type:      types.CounterKind,
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
