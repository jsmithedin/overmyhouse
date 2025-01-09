package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type aircraftData struct {
	icaoAddr uint32

	callsign string

	eRawLat uint32
	eRawLon uint32
	oRawLat uint32
	oRawLon uint32

	latitude  float64
	longitude float64
	altitude  int32

	lastPing time.Time
	lastPos  time.Time

	mlat bool
}

type aircraftList []*aircraftData
type aircraftMap map[uint32]*aircraftData

type KnownAircraft struct {
	knownMap aircraftMap
	mu       sync.Mutex
}

func (kAircraft *KnownAircraft) getNumberOfKnown() (total int) {
	kAircraft.mu.Lock()
	defer kAircraft.mu.Unlock()
	total = len(kAircraft.knownMap)
	return total
}

func (kAircraft *KnownAircraft) getAircraft(icaoAddr uint32) (ptrAircraft *aircraftData, aircraftExists bool) {
	kAircraft.mu.Lock()
	defer kAircraft.mu.Unlock()
	ptrAircraft, aircraftExists = kAircraft.knownMap[icaoAddr]
	return ptrAircraft, aircraftExists
}

func (kAircraft *KnownAircraft) addAircraft(icaoAddr uint32, aircraft *aircraftData) {
	kAircraft.mu.Lock()
	if kAircraft.knownMap == nil {
		kAircraft.knownMap = make(aircraftMap)
	}

	kAircraft.knownMap[icaoAddr] = aircraft
	kAircraft.mu.Unlock()
}

func (kAircraft *KnownAircraft) removeAircraft(icaoAddr uint32) {
	kAircraft.mu.Lock()
	delete(kAircraft.knownMap, icaoAddr)
	kAircraft.mu.Unlock()
}

func (kAircraft *KnownAircraft) pruneKnown(now time.Time, timeout uint32) {
	knownList := kAircraft.sortedAircraft()
	for _, aircraft := range knownList {
		if now.Sub(aircraft.lastPing).Seconds() > float64(timeout) {
			kAircraft.removeAircraft(aircraft.icaoAddr)
		}
	}
}

func (kAircraft *KnownAircraft) sortedAircraft() (sortedAircraftList aircraftList) {
	kAircraft.mu.Lock()
	sortedAircraftList = make(aircraftList, 0, len(kAircraft.knownMap))

	for _, aircraft := range kAircraft.knownMap {
		sortedAircraftList = append(sortedAircraftList, aircraft)
	}

	sort.Sort(sortedAircraftList)
	kAircraft.mu.Unlock()
	return sortedAircraftList
}

func (a aircraftList) Len() int {
	return len(a)
}

func (a aircraftList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a aircraftList) Less(i, j int) bool {
	if a[i].latitude != math.MaxFloat64 && a[j].latitude != math.MaxFloat64 {
		return sortAircraftByDistance(a, i, j)
	} else if a[i].latitude != math.MaxFloat64 && a[j].latitude == math.MaxFloat64 {
		return true
	} else if a[i].latitude == math.MaxFloat64 && a[j].latitude != math.MaxFloat64 {
		return false
	}
	return sortAircraftByCallsign(a, i, j)
}

func sortAircraftByDistance(a aircraftList, i, j int) bool {
	return GreatCircle(a[i].latitude, a[i].longitude, *baseLat, *baseLon) <
		GreatCircle(a[j].latitude, a[j].longitude, *baseLat, *baseLon)
}

func sortAircraftByCallsign(a aircraftList, i, j int) bool {
	if a[i].callsign != "" && a[j].callsign != "" {
		return a[i].callsign < a[j].callsign
	} else if a[i].callsign != "" && a[j].callsign == "" {
		return true
	} else if a[i].callsign == "" && a[j].callsign != "" {
		return false
	}
	return fmt.Sprintf("%06x", a[i].icaoAddr) < fmt.Sprintf("%06x", a[j].icaoAddr)
}
