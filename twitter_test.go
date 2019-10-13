package main

import (
	"testing"
	"time"
)

var testTweeted *TweetedAircraft

func TestAddingAircraft(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.AddAircraft("addAircraft")
	if testTweeted.tweetedMap["addAircraft"] == 0 {
		t.Errorf("Didn't add aircraft")
	}
}

func TestCheckForAircraft(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.AddAircraft("checkAircraft")

	if !testTweeted.AlreadyTweeted("checkAircraft") {
		t.Errorf("Didn't properly find aircraft")
	}

	if testTweeted.AlreadyTweeted("checkAnotherAircraft") {
		t.Errorf("Found an aircraft I shouldn't have")
	}
}

func TestPruneAircraft(t *testing.T) {
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

func TestConcurrentAdd(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.AddAircraft("ABC")
	channel := make(chan bool)

	for i := 0; i < 5; i++ {
		go add(testTweeted, "ABC"+string(i), channel)
	}

	channel <- true

	if !testTweeted.AlreadyTweeted("ABC") {
		t.Errorf("Concurrent adding went wrong")
	}
}

func TestConcurrentAddAndPrune(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.AddAircraft("ABC")
	channel := make(chan bool)

	for i := 0; i < 5; i++ {
		go add(testTweeted, "ABC"+string(i), channel)
	}

	go prune(testTweeted, channel)

	channel <- true

	if !testTweeted.AlreadyTweeted("ABC") {
		t.Errorf("Concurrent adding went wrong")
	}
}

func add(tt *TweetedAircraft, cs string, c <-chan bool) {
	if <-c {
		tt.AddAircraft(cs)
	}
}

func prune(tt *TweetedAircraft, c <-chan bool) {
	if <-c {
		tt.PruneTweeted()
	}
}
