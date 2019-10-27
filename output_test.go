package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_durationSecondsElapse(t *testing.T) {
	tests := []struct {
		since time.Duration
		want  string
	}{
		{since: time.Duration(255000000000), want: "-"},
		{since: time.Duration(1000000000), want: "   1"},
	}

	for _, tc := range tests {
		got := durationSecondsElapsed(tc.since)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v, secs: %d", tc.want, got, uint8(tc.since.Seconds()))
		}
	}
}
