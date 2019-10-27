package main

import (
	"math"
	"sync"
	"testing"
)

var testKnown *KnownAircraft

func TestAddingAircraftToKnown(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}
	testKnown.addAircraft(123, &testAircraft)
	if testKnown.knownMap[123] == nil {
		t.Errorf("Aircraft not added to known")
	}
}

func TestGetNumberOfKnown(t *testing.T) {
	testKnown = &KnownAircraft{}
	testAircraft := aircraftData{}

	for i := 0; i < 5; i++ {
		testKnown.addAircraft(uint32(i), &testAircraft)
	}

	total := testKnown.getNumberOfKnown()

	if total != 5 {
		t.Errorf("Incorrect total of aircraft known got %d", total)
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

	for i := 0; i < 5; i++ {
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

	wg := sync.WaitGroup{}

	wg.Add(1)
	go addKnown(testKnown, 123, &testAircraft, &wg)
	wg.Add(1)
	go sortKnown(testKnown, &wg)

	wg.Wait()
}

func addKnown(ka *KnownAircraft, icao uint32, ac *aircraftData, wg *sync.WaitGroup) {
	ka.addAircraft(icao, ac)
	wg.Done()
}

func sortKnown(ka *KnownAircraft, wg *sync.WaitGroup) {
	_ = ka.sortedAircraft()
	wg.Done()
}

func setupAircraftList() aircraftList {
	testList := make(aircraftList, 0, 10)
	return testList
}

func TestSwap(t *testing.T) {
	testList := setupAircraftList()
	testAircraft1 := aircraftData{icaoAddr: 1, callsign: "b", latitude: math.MaxFloat64}
	testAircraft2 := aircraftData{icaoAddr: 2, callsign: "a", latitude: math.MaxFloat64}

	testList = append(testList, &testAircraft1)
	testList = append(testList, &testAircraft2)

	testList.Swap(0, 1)

	if testList[0].icaoAddr != 2 {
		t.Errorf("Didn't properly swap")
	}
}

func TestLessNoLocation(t *testing.T) {
	testList := setupAircraftList()
	testAircraft1 := aircraftData{icaoAddr: 1, callsign: "b", latitude: math.MaxFloat64}
	testAircraft2 := aircraftData{icaoAddr: 2, callsign: "a", latitude: math.MaxFloat64}

	testList = append(testList, &testAircraft1)
	testList = append(testList, &testAircraft2)

	less := testList.Less(0, 1)

	if less != false {
		t.Errorf("Didn't less")
	}
}

func TestLessOneLocation(t *testing.T) {
	testList := setupAircraftList()
	testAircraft1 := aircraftData{icaoAddr: 2, callsign: "a", latitude: math.MaxFloat64}
	testAircraft2 := aircraftData{icaoAddr: 2, callsign: "c", latitude: 1.0}
	testList = append(testList, &testAircraft1)
	testList = append(testList, &testAircraft2)

	less := testList.Less(1, 0)

	if less != true {
		t.Errorf("Didn't properly less")
	}

	less = testList.Less(0, 1)

	if less != false {
		t.Errorf("Didn't properly less")
	}
}

func TestLessTwoLocations(t *testing.T) {
	testList := setupAircraftList()
	testAircraft1 := aircraftData{icaoAddr: 2, callsign: "c", latitude: 1.0, longitude: 1.0}
	testAircraft2 := aircraftData{icaoAddr: 2, callsign: "d", latitude: 2.0, longitude: 2.0}
	testList = append(testList, &testAircraft1)
	testList = append(testList, &testAircraft2)

	less := testList.Less(0, 1)

	if less != false {
		t.Errorf("Didn't properly less")
	}
}

func TestLessByCallsign(t *testing.T) {
	testList := setupAircraftList()
	testAircraft1 := aircraftData{icaoAddr: 1, callsign: "c", latitude: math.MaxFloat64}
	testAircraft2 := aircraftData{icaoAddr: 2, callsign: "d", latitude: math.MaxFloat64}
	testAircraft3 := aircraftData{icaoAddr: 3, latitude: math.MaxFloat64}
	testAircraft4 := aircraftData{icaoAddr: 4, latitude: math.MaxFloat64}
	testList = append(testList, &testAircraft1)
	testList = append(testList, &testAircraft2)
	testList = append(testList, &testAircraft3)
	testList = append(testList, &testAircraft4)

	less := testList.Less(0, 1)

	if less != true {
		t.Errorf("Didn't properly less")
	}

	less = testList.Less(2, 0)

	if less != false {
		t.Errorf("Didn't properly less")
	}

	less = testList.Less(0, 2)

	if less != true {
		t.Errorf("Didn't properly less")
	}

	less = testList.Less(3, 2)

	if less != false {
		t.Errorf("Didn't properly less")
	}
}
