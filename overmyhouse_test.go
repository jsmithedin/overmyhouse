package main

import (
	"testing"
	"time"
)

// Test that startClient does not panic when connection cannot be established.
func TestStartClient_NoPanic(t *testing.T) {
	ch := startClient("invalid:0")
	select {
	case conn := <-ch:
		if conn != nil {
			conn.Close()
		}
	case <-time.After(1 * time.Second):
		t.Fatal("startClient did not return")
	}
}
