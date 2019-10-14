package main

import (
	"fmt"
	"sync"
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
	channel := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go addTweeted(testTweeted, channel, &wg)
	}

	for i := 0; i < 500; i++ {
		channel <- fmt.Sprintf("ABC%d", i)
	}

	wg.Wait()

	for i := 0; i < 500; i++ {
		if !testTweeted.AlreadyTweeted(fmt.Sprintf("ABC%d", i)) {
			t.Errorf("Concurrent adding went wrong")
		}
	}
}

func TestConcurrentAddAndPruneTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.AddAircraft("ABC")
	channel := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go addTweeted(testTweeted, channel, &wg)
	}

	wg.Add(1)
	go pruneTweeted(testTweeted, &wg)

	for i := 0; i < 500; i++ {
		channel <- fmt.Sprintf("ABC%d", i)
	}

	wg.Wait()

	for i := 0; i < 500; i++ {
		if !testTweeted.AlreadyTweeted(fmt.Sprintf("ABC%d", i)) {
			t.Errorf("Concurrent adding went wrong")
		}
	}
}

func addTweeted(tt *TweetedAircraft, c <-chan string, wg *sync.WaitGroup) {
	cs := <-c
	tt.AddAircraft(cs)
	wg.Done()
}

func pruneTweeted(tt *TweetedAircraft, wg *sync.WaitGroup) {
	tt.PruneTweeted()
	wg.Done()
}
