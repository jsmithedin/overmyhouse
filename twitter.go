package main

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

type tweetedMap map[string]int64

// TweetedAircraft ties a map of Aircraft we have already tweeted about with a mutex controlling access to the map
type TweetedAircraft struct {
	tweetedMap tweetedMap
	mu         sync.Mutex
}

// AddAircraft adds a new Aircraft to the map of Aircraft we have tweeted about
func (tAircraft *TweetedAircraft) AddAircraft(callsign string) {
	tAircraft.mu.Lock()

	if tAircraft.tweetedMap == nil {
		tAircraft.tweetedMap = make(tweetedMap)
	}

	tAircraft.tweetedMap[callsign] = time.Now().Unix()

	tAircraft.mu.Unlock()
}

// AlreadyTweeted checks if an aircraft is present in the map
func (tAircraft *TweetedAircraft) AlreadyTweeted(callsign string) bool {
	tAircraft.mu.Lock()
	defer tAircraft.mu.Unlock()
	_, ok := tAircraft.tweetedMap[callsign]
	return ok
}

// PruneTweeted removes aircraft from the map which have been present for 60s
func (tAircraft *TweetedAircraft) PruneTweeted() {
	tAircraft.mu.Lock()
	timeNow := time.Now().Unix()

	for callsign, timeAdded := range (*tAircraft).tweetedMap {
		if (timeNow - timeAdded) > 60 {
			delete(tAircraft.tweetedMap, callsign)
		}
	}

	tAircraft.mu.Unlock()
}

func tweet(message string) (int64, error) {
	err := godotenv.Load()
	if err != nil {
		return 0, err
	}

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
	}

	return tweet.ID, nil
}
