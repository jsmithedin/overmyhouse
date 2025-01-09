package main

import (
	"encoding/binary"
	"math"
	"reflect"
	"testing"
	"time"
)

func Test_ParseTime(t *testing.T) {
	const testTimestampValue = 0x244bbb9ac9f0 // Example GPS timestamp in bytes
	expectedUnixTimestamp := int64(1676728358)

	// Extracted helper to create timestamp bytes
	timestampBytes := createTimestampBytes(testTimestampValue)

	utcDate := time.Unix(expectedUnixTimestamp, 0).UTC()
	timestamp := parseTime(timestampBytes, utcDate)

	if timestamp.Unix() != expectedUnixTimestamp {
		t.Errorf("Expected %d, but got %d", expectedUnixTimestamp, timestamp.Unix())
	}
}

// Helper function to create timestamp byte slice from uint64 value
func createTimestampBytes(value uint64) []byte {
	timestampBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestampBytes, value)
	return timestampBytes
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
		{message: []byte{141, 64, 115, 119, 232, 52, 66, 112, 226, 8, 32, 34, 235, 246}, callsign: "MDI08   ", altitude: 6500, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 40, 52, 66, 112, 226, 8, 32, 34, 235, 246}, callsign: "MDI08   ", altitude: 6500, latitude: 219.6614227294922, longitude: 5.712890625},
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

func Test_setPositions(t *testing.T) {
	testAircraft := aircraftData{}
	testAircrafteLatLon := aircraftData{eRawLat: math.MaxUint32, eRawLon: math.MaxUint32}
	testAircraftoLatLon := aircraftData{oRawLat: math.MaxUint32, oRawLon: math.MaxUint32}

	tests := []struct {
		message      []byte
		testAircraft aircraftData
		rawLatitude  uint32
		rawLongitude uint32
		latitude     float64
		longitude    float64
	}{
		{message: []byte{141, 64, 115, 119, 33, 52, 66, 112, 226, 8, 32, 34, 235, 246}, testAircraft: testAircraft, rawLatitude: 0, rawLongitude: 0, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 33, 52, 60, 112, 226, 8, 32, 34, 235, 246}, testAircraft: testAircraft, rawLatitude: 0, rawLongitude: 0, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 33, 52, 66, 112, 226, 8, 32, 34, 235, 246}, testAircraft: testAircrafteLatLon, rawLatitude: 0, rawLongitude: 0, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 33, 52, 60, 112, 226, 8, 32, 34, 235, 246}, testAircraft: testAircrafteLatLon, rawLatitude: 0, rawLongitude: 0, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 33, 52, 66, 112, 226, 8, 32, 34, 235, 246}, testAircraft: testAircraftoLatLon, rawLatitude: 0, rawLongitude: 0, latitude: 0, longitude: 0},
		{message: []byte{141, 64, 115, 119, 33, 52, 60, 112, 226, 8, 32, 34, 235, 246}, testAircraft: testAircraftoLatLon, rawLatitude: 0, rawLongitude: 0, latitude: 0, longitude: 0},
	}

	for _, tc := range tests {
		latGot, lonGot := setPositions(&tc.message, &tc.testAircraft, tc.rawLatitude, tc.rawLongitude)
		if !reflect.DeepEqual(latGot, tc.latitude) {
			t.Fatalf("expected: %v, got: %v", tc.latitude, latGot)
		}
		if !reflect.DeepEqual(lonGot, tc.longitude) {
			t.Fatalf("expected: %v, got %v", tc.longitude, lonGot)
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
		{message: []byte{0, 64, 15, 154, 153, 20, 254, 133, 161, 36, 130, 240, 148, 109}, isMlat: true, number: 1},
	}

	for _, tc := range tests {
		parseModeS(tc.message, tc.isMlat, testKnownAircraft)
		if !reflect.DeepEqual(testKnownAircraft.getNumberOfKnown(), tc.number) {
			t.Fatalf("expected: %v, got: :%v:", tc.number, testKnownAircraft.getNumberOfKnown())
		}
	}
}
