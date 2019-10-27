package main

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"
)

func Test_ParseTime(t *testing.T) {
	timestampBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestampBytes, 0x244bbb9ac9f0)
	timestamp := parseTime(timestampBytes)

	if timestamp.Unix() != 1572184358 {
		t.Errorf("Got %d", timestamp.Unix())
	}
}

func Test_ParseRawLatLon(t *testing.T) {
	tests := []struct {
		evenLat uint32
		evenLon uint32
		oddLat  uint32
		oddLon  uint32
		lastOdd bool
		tFlag   bool
		latWant float64
		lonWant float64
	}{
		{evenLat: 92095, evenLon: 39846, oddLat: 88385, oddLon: 125818, lastOdd: false, tFlag: false, latWant: 10.215774536132812, lonWant: 123.88881877317269},
		{evenLat: 92095, evenLon: 39846, oddLat: 88385, oddLon: 125818, lastOdd: false, tFlag: true, latWant: 10.21621445478019, lonWant: 123.8891285863416},
		{evenLat: 92095, evenLon: 39846, oddLat: 88385, oddLon: 125818, lastOdd: true, tFlag: false, latWant: 10.215774536132812, lonWant: 123.88881877317269},
		{evenLat: 92095, evenLon: 39846, oddLat: 88385, oddLon: 125818, lastOdd: true, tFlag: true, latWant: 10.21621445478019, lonWant: 123.8891285863416},
		{evenLat: 92095, evenLon: 39846, oddLat: 88385, oddLon: math.MaxUint32, lastOdd: false, tFlag: false, latWant: math.MaxFloat64, lonWant: math.MaxFloat64},
	}

	for _, tc := range tests {
		latGot, lonGot := parseRawLatLon(tc.evenLat, tc.evenLon, tc.oddLat, tc.oddLon, tc.lastOdd, tc.tFlag)
		if !reflect.DeepEqual(tc.latWant, latGot) {
			t.Fatalf("expected: %v, got: %v", tc.latWant, latGot)
		}
		if !reflect.DeepEqual(tc.lonWant, lonGot) {
			t.Fatalf("expected: %v, got: %v", tc.lonWant, lonGot)
		}
	}
}
