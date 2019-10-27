package main

import (
	"encoding/binary"
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
