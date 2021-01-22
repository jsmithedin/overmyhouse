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
	testTweeted.addAircraft("addAircraft")
	if testTweeted.tweetedMap["addAircraft"] == 0 {
		t.Errorf("Didn't add aircraft")
	}
}

func TestGetNumberOfTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	for i := 0; i < 5; i++ {
		testTweeted.addAircraft(string(rune(i)))
	}

	total := testTweeted.getNumberOfTweeted()
	if total != 5 {
		t.Errorf("Incorrect total of aircraft tweeted got %d", total)
	}
}

func TestCheckForAircraftTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.addAircraft("checkAircraft")

	if !testTweeted.alreadyTweeted("checkAircraft") {
		t.Errorf("Didn't properly find aircraft")
	}

	if testTweeted.alreadyTweeted("checkAnotherAircraft") {
		t.Errorf("Found an aircraft I shouldn't have")
	}
}

func TestPruneAircraftTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	testTweeted.addAircraft("pruneAircraft")
	testTweeted.addAircraft("dontPruneAircraft")
	testTweeted.tweetedMap["pruneAircraft"] = time.Now().Unix() - 120

	testTweeted.pruneTweeted()

	if testTweeted.alreadyTweeted("pruneAircraft") {
		t.Errorf("Didn't properly prune")
	}

	if !testTweeted.alreadyTweeted("dontPruneAircraft") {
		t.Errorf("Pruned too many")
	}
}

func TestConcurrentAddTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.addAircraft("ABC")
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
		if !testTweeted.alreadyTweeted(fmt.Sprintf("ABC%d", i)) {
			t.Errorf("Concurrent adding went wrong")
		}
	}
}

func TestConcurrentAddAndPruneTweeted(t *testing.T) {
	testTweeted = &TweetedAircraft{}
	// Start with one aircraft to initialise
	testTweeted.addAircraft("ABC")
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
		if !testTweeted.alreadyTweeted(fmt.Sprintf("ABC%d", i)) {
			t.Errorf("Concurrent adding went wrong")
		}
	}
}

func addTweeted(tt *TweetedAircraft, c <-chan string, wg *sync.WaitGroup) {
	cs := <-c
	tt.addAircraft(cs)
	wg.Done()
}

func pruneTweeted(tt *TweetedAircraft, wg *sync.WaitGroup) {
	tt.pruneTweeted()
	wg.Done()
}
