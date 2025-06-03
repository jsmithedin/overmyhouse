package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSlackNotify(t *testing.T) {
	var received struct {
		Text string `json:"text"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		_ = json.Unmarshal(body, &received)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	os.Setenv("slackwebhook", server.URL)
	err := slackNotify("test message")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if received.Text != "test message" {
		t.Fatalf("expected 'test message', got %s", received.Text)
	}
}

func TestSlackNotifyMissingEnv(t *testing.T) {
	os.Unsetenv("slackwebhook")
	err := slackNotify("msg")
	if err == nil {
		t.Fatalf("expected error when env missing")
	}
}
