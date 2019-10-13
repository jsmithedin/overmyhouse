package main

import (
	"testing"
	"time"
)

var testTweeted *TweetedAircraft

func TestAddingAircraftToTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.AddAircraft("addAircraft")
	if testTweeted.tweetedMap["addAircraft"] == 0 {
		t.Errorf("Didn't add aircraft")
	}
}

func TestCheckForAircraftTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.AddAircraft("checkAircraft")

	if !testTweeted.AlreadyTweeted("checkAircraft") {
		t.Errorf("Didn't properly find aircraft")
	}

	if testTweeted.AlreadyTweeted("checkAnotherAircraft") {
		t.Errorf("Found an aircraft I shouldn't have")
	}
}

func TestPruneAircraftTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.AddAircraft("pruneAircraft")
	testTweeted.AddAircraft("dontPruneAircraft")
	testTweeted.tweetedMap["pruneAircraft"] = time.Now().Unix() - 120

	testTweeted.PruneTweeted()

	if testTweeted.AlreadyTweeted("pruneAircraft") {
		t.Errorf("Didn't properly prune")
	}

	if !testTweeted.AlreadyTweeted("dontPruneAircraft") {
		t.Errorf("Pruned too many")
	}
}

func TestConcurrentAddTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.AddAircraft("ABC")
	channel := make(chan bool)

	for i := 0; i < 5; i++ {
		go addTweeted(testTweeted, "ABC"+string(i), channel)
	}

	channel <- true

	if !testTweeted.AlreadyTweeted("ABC") {
		t.Errorf("Concurrent adding went wrong")
	}
}

func TestConcurrentAddAndPruneTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.AddAircraft("ABC")
	channel := make(chan bool)

	for i := 0; i < 5; i++ {
		go addTweeted(testTweeted, "ABC"+string(i), channel)
	}

	go pruneTweeted(testTweeted, channel)

	channel <- true

	if !testTweeted.AlreadyTweeted("ABC") {
		t.Errorf("Concurrent adding went wrong")
	}
}

func addTweeted(tt *TweetedAircraft, cs string, c <-chan bool) {
	if <-c {
		tt.AddAircraft(cs)
	}
}

func pruneTweeted(tt *TweetedAircraft, c <-chan bool) {
	if <-c {
		tt.PruneTweeted()
	}
}
