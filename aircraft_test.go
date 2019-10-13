package main

import "testing"

var testKnown *KnownAircraft

func TestAddingAircraftToKnown(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}
	testKnown.addAircraft(123, &testAircraft)
	if testKnown.knownMap[123] == nil {
		t.Errorf("Aircraft not added to known")
	}
}

func TestRemoveAircraftFromKnown(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}
	testKnown.addAircraft(123, &testAircraft)

	testKnown.removeAircraft(123)

	if testKnown.knownMap[123] != nil {
		t.Errorf("Aircraft not removed from known")
	}
}

func TestGetAircraftKnown(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}
	testKnown.addAircraft(123, &testAircraft)

	returnedAircraft, known := testKnown.getAircraft(123)
	if returnedAircraft == nil || known == false {
		t.Errorf("Couldn't get known aircraft")
	}

	returnedAircraft, known = testKnown.getAircraft(456)
	if returnedAircraft != nil || known == true {
		t.Errorf("Got an unknown aircraft")
	}
}

func TestGetSortedAircraft(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}

	for i := 5; i < 0; i-- {
		testKnown.addAircraft(uint32(i), &testAircraft)
	}

	sorted := testKnown.sortedAircraft()

	if sorted == nil {
		t.Errorf("Didn't get sorted aircraft")
	}
}

func TestSortWhilstAdding(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}
	channel := make(chan bool)

	go addKnown(testKnown, 123, &testAircraft, channel)
	go sortKnown(testKnown, channel)

	channel <- true
}

func addKnown(ka *KnownAircraft, icao uint32, ac *aircraftData, c <-chan bool) {
	if <-c {
		ka.addAircraft(icao, ac)
	}
}

func sortKnown(ka *KnownAircraft, c <-chan bool) {
	if <-c {
		_ = ka.sortedAircraft()
	}
}
