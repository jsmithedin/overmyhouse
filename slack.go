package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// slackNotify posts a message to a Slack incoming webhook.
// The webhook URL is read from the `slackwebhook` environment variable.
func slackNotify(message string) error {
	_ = godotenv.Load()

	webhookURL := os.Getenv("slackwebhook")
	if len(webhookURL) == 0 {
		return errors.New("Env isnt set correctly")
	}

	payload := map[string]string{"text": message}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("failed to post message to slack")
	}

	return nil
}
