package main

import (
	"reflect"
	"testing"
)

func TestKV(t *testing.T) {
	testCases := []struct {
		in  []string
		out map[string]string
	}{
		{
			[]string{"key=value"},
			map[string]string{"key": "value"},
		},
		{
			[]string{"k1=v1", "k2=v2"},
			map[string]string{"k1": "v1", "k2": "v2"},
		},
		{
			[]string{"key=k=v"},
			map[string]string{"key": "k=v"},
		},
	}

	for _, tc := range testCases {
		out := make(map[string]string)
		kv := NewKV(out)
		for _, in := range tc.in {
			if err := kv.Set(in); err != nil {
				t.Errorf("fail to set %s: %n", in, err)
			}
		}

		if !reflect.DeepEqual(out, tc.out) {
			t.Errorf("fail to set %v.\nGot: %v\nWant: %v", tc.in, out, tc.out)
		}
	}
}
