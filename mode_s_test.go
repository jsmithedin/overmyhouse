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

	// TODO: Investigate why this changed once...
	if timestamp.Unix() != 1574862758 {
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

func Test_decodeExtendedSquitter(t *testing.T) {
	testAircraft := aircraftData{}

	tests := []struct {
		message   []byte
		callsign  string
		altitude  int32
		latitude  float64
		longitude float64
	}{
		{message: []byte{141, 64, 12, 74, 153, 68, 2, 22, 232, 72, 11, 144, 151, 165}, callsign: "", altitude: 0, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 86, 11, 88, 37, 196, 163, 243, 90, 151, 218, 105, 13}, callsign: "", altitude: 6500, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 33, 52, 66, 112, 226, 8, 32, 34, 235, 246}, callsign: "MDI08   ", altitude: 6500, latitude: 0, longitude: 0},
	}

	for _, tc := range tests {
		decodeExtendedSquitter(tc.message, &testAircraft)
		if !reflect.DeepEqual(testAircraft.callsign, tc.callsign) {
			t.Fatalf("expected: %v, got: :%v:", tc.callsign, testAircraft.callsign)
		}
		if !reflect.DeepEqual(testAircraft.altitude, tc.altitude) {
			t.Fatalf("expected: %v, got: %v", tc.altitude, testAircraft.altitude)
		}
		if !reflect.DeepEqual(testAircraft.latitude, tc.latitude) {
			t.Fatalf("expected: %v, got: %v", tc.latitude, testAircraft.latitude)
		}
		if !reflect.DeepEqual(testAircraft.longitude, tc.longitude) {
			t.Fatalf("expected: %v, got: %v", tc.longitude, testAircraft.longitude)
		}
	}
}

func Test_parseModeS(t *testing.T) {
	testKnownAircraft := &KnownAircraft{}

	tests := []struct {
		message []byte
		isMlat  bool
		number  int
	}{
		{message: []byte{141, 64, 15, 154, 153, 20, 254, 133, 161, 36, 130, 240, 148, 109}, isMlat: true, number: 1},
		{message: []byte{141, 64, 15, 154, 153, 20, 254, 133, 161, 36, 130, 240, 148, 109}, isMlat: true, number: 1},
	}

	for _, tc := range tests {
		parseModeS(tc.message, tc.isMlat, testKnownAircraft)
		if !reflect.DeepEqual(testKnownAircraft.getNumberOfKnown(), tc.number) {
			t.Fatalf("expected: %v, got: :%v:", tc.number, testKnownAircraft.getNumberOfKnown())
		}
	}
}
