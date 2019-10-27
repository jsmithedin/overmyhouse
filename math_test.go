package main

import (
	"reflect"
	"testing"
)

func Test_cprNLFunction(t *testing.T) {
	tests := []struct {
		input float64
		want  byte
	}{
		{input: -1.0, want: 59},
		{input: 12.0, want: 58},
		{input: 16.0, want: 57},
		{input: 19.0, want: 56},
		{input: 22.0, want: 55},
		{input: 24.0, want: 54},
		{input: 26.0, want: 53},
		{input: 28.0, want: 52},
		{input: 30.0, want: 51},
		{input: 32.0, want: 50},
		{input: 34.0, want: 49},
		{input: 36.0, want: 48},
		{input: 37.0, want: 47},
		{input: 39.0, want: 46},
		{input: 40.0, want: 45},
		{input: 42.0, want: 44},
		{input: 43.0, want: 43},
		{input: 45.0, want: 42},
		{input: 46.0, want: 41},
		{input: 47.0, want: 40},
		{input: 49.0, want: 39},
		{input: 50.0, want: 38},
		{input: 51.0, want: 37},
		{input: 53.0, want: 36},
		{input: 54.0, want: 35},
		{input: 55.0, want: 34},
		{input: 56.0, want: 33},
		{input: 57.0, want: 32},
		{input: 58.0, want: 31},
		{input: 59.0, want: 30},
		{input: 61.0, want: 29},
		{input: 62.0, want: 28},
		{input: 63.0, want: 27},
		{input: 64.0, want: 26},
		{input: 65.0, want: 25},
		{input: 66.0, want: 24},
		{input: 67.0, want: 23},
		{input: 68.0, want: 22},
		{input: 69.0, want: 21},
		{input: 70.0, want: 20},
		{input: 71.0, want: 19},
		{input: 72.0, want: 18},
		{input: 73.0, want: 17},
		{input: 74.0, want: 16},
		{input: 75.0, want: 15},
		{input: 76.0, want: 14},
		{input: 77.0, want: 13},
		{input: 78.0, want: 12},
		{input: 79.0, want: 11},
		{input: 80.0, want: 10},
		{input: 81.0, want: 9},
		{input: 82.0, want: 8},
		{input: 83.0, want: 7},
		{input: 83.5, want: 6},
		{input: 84.5, want: 5},
		{input: 85.5, want: 4},
		{input: 86.5, want: 3},
		{input: 86.6, want: 2},
		{input: 89.0, want: 1},
	}

	for _, tc := range tests {
		got := cprNLFunction(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}
}

func Test_cprNFunction(t *testing.T) {
	tests := []struct {
		input float64
		fflag bool
		want  byte
	}{
		{input: 5.0, fflag: false, want: 59},
		{input: 5.0, fflag: true, want: 58},
		{input: 88.0, fflag: true, want: 1},
	}

	for _, tc := range tests {
		got := cprNFunction(tc.input, tc.fflag)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}

}

func Test_cprDlonFunction(t *testing.T) {
	tests := []struct {
		input   float64
		fflag   bool
		surface bool
		want    float64
	}{
		{input: 5.0, fflag: true, surface: true, want: 1.5517241379310345},
		{input: 5.0, fflag: true, surface: false, want: 6.206896551724138},
	}

	for _, tc := range tests {
		got := cprDlonFunction(tc.input, tc.fflag, tc.surface)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}

}

func Test_decodeAC12Field(t *testing.T) {
	tests := []struct {
		input uint
		want  int32
	}{
		{input: 1, want: 2147483647},
		{input: 16, want: -1000},
	}

	for _, tc := range tests {
		got := decodeAC12Field(tc.input)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}

}

func Test_greatcircle(t *testing.T) {
	tests := []struct {
		lat0 float64
		lon0 float64
		lat1 float64
		lon1 float64
		want float64
	}{
		{lat0: 1.0, lon0: 1.0, lat1: 1.0, lon1: 1.0, want: 0},
		{lat0: 1.0, lon0: 1.0, lat1: 2.0, lon1: 2.0, want: 157225.43203804837},
	}

	for _, tc := range tests {
		got := greatcircle(tc.lat0, tc.lon0, tc.lat1, tc.lon1)
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("expected: %v, got: %v", tc.want, got)
		}
	}

}

func Test_MetersInMiles(t *testing.T) {
	miles := metersInMiles(1609.34721869)

	if miles != 1 {
		t.Errorf("Incorrect meters in miles Got %f", miles)
	}
}
