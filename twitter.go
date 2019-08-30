package main

import (
	"errors"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func tweet(message string) (int64, error) {
	consumerKey := os.Getenv("consumerkey")
	consumerSecret := os.Getenv("consumersecret")
	accessToken := os.Getenv("accesstoken")
	accessSecret := os.Getenv("accesssecret")

	if len(consumerKey) == 0 || len(consumerSecret) == 0 || len(accessToken) == 0 || len(accessSecret) == 0 {
		return 0, errors.New("Env isnt set correctly")
	}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	tweet, _, err := client.Statuses.Update(message, nil)

	if err != nil {
		return 0, err
	} else {
		return tweet.ID, nil
	}
}

