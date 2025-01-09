package main

import (
	"fmt"
	"log"
	"math"
	"time"
)

func durationSecondsElapsed(since time.Duration) string {
	sec := uint8(since.Seconds())
	if sec == math.MaxUint8 {
		return "-"
	} else {
		return fmt.Sprintf("%4d", sec)
	}
}

func printStats(knownAircraft *KnownAircraft, tweetedAircraft *TweetedAircraft) {
	t := time.Now()
	numberOfKnownAircraft := knownAircraft.getNumberOfKnown()
	numberOfTweetedAircraft := tweetedAircraft.getNumberOfTweeted()

	fmt.Printf("%d-%02d-%02dT%02d:%02d:%02d-00:00 Known: %d\tTweeted: %d\n", t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), numberOfKnownAircraft, numberOfTweetedAircraft)
}

func printOverhead(knownAircraft *KnownAircraft, tweetedAircraft *TweetedAircraft, radius *int) {
	sortedAircraft := knownAircraft.sortedAircraft()

	for icao, aircraft := range sortedAircraft {
		stale := (time.Since(aircraft.lastPos) > time.Duration((10)*time.Second))
		extraStale := (time.Since(aircraft.lastPos) > (time.Duration(20) * time.Second))

		aircraftHasLocation := (aircraft.latitude != math.MaxFloat64 &&
			aircraft.longitude != math.MaxFloat64)
		aircraftHasAltitude := aircraft.altitude != math.MaxInt32

		if aircraftHasLocation && aircraftHasAltitude {
			sLatLon := fmt.Sprintf("%f,%f", aircraft.latitude, aircraft.longitude)
			sAlt := fmt.Sprintf("%d", aircraft.altitude)

			distance := GreatCircle(aircraft.latitude, aircraft.longitude,
				*baseLat, *baseLon)

			tPos := time.Since(aircraft.lastPos)

			if !stale && !extraStale && metersInMiles(distance) < float64(*radius) {
				if !tweetedAircraft.alreadyTweeted(aircraft.callsign) {
					log.Printf("%06x\t%8s\t%s%s\t%3.2f\t%s\n",
						aircraft.icaoAddr, aircraft.callsign,
						sLatLon, sAlt, metersInMiles(distance),
						durationSecondsElapsed(tPos))

					if len(aircraft.callsign) > 0 {
						_, err := tweet(fmt.Sprintf("https://flightaware.com/live/flight/%8s %8s flew %3.2f miles from my house at %d ft!",
							aircraft.callsign, aircraft.callsign, metersInMiles(distance), aircraft.altitude))

						if err != nil {
							log.Print(err)
						}

						tweetedAircraft.addAircraft(aircraft.callsign)
					}
				}
			}
			if extraStale {
				knownAircraft.removeAircraft(uint32(icao))
			}
		}
	}
}

func printAircraftTable(knownAircraft *KnownAircraft) {
	fmt.Print("\x1b[H\x1b[2J")
	fmt.Println("ICAO \tCallsign\tLocation\t\tAlt\tDistance   Time")

	sortedAircraft := knownAircraft.sortedAircraft()

	for _, aircraft := range sortedAircraft {
		stale := (time.Since(aircraft.lastPos) > time.Duration((10)*time.Second))
		extraStale := (time.Since(aircraft.lastPos) > (time.Duration(20) * time.Second))

		aircraftHasLocation := (aircraft.latitude != math.MaxFloat64 &&
			aircraft.longitude != math.MaxFloat64)
		aircraftHasAltitude := aircraft.altitude != math.MaxInt32

		if aircraft.callsign != "" || aircraftHasLocation || aircraftHasAltitude {
			var sLatLon string
			var sAlt string

			if aircraftHasLocation {
				sLatLon = fmt.Sprintf("%f,%f", aircraft.latitude, aircraft.longitude)
			} else {
				sLatLon = "---.------, ---.------"
			}
			if aircraftHasAltitude {
				sAlt = fmt.Sprintf("%d", aircraft.altitude)
			} else {
				sAlt = "-----"
			}

			distance := GreatCircle(aircraft.latitude, aircraft.longitude,
				*baseLat, *baseLon)

			isMlat := ""
			if aircraft.mlat {
				isMlat = "^"
			}

			tPos := time.Since(aircraft.lastPos)

			if !stale && !extraStale {
				fmt.Printf("%06x\t%8s\t%s%s\t%s\t%3.2f\t%s\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, isMlat, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			} else if stale && !extraStale {
				fmt.Printf("%06x\t%8s\t%s%s?\t%s\t%3.2f?\t%s\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, isMlat, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			} else {
				fmt.Printf("%06x\t%8s\t%s%s?\t%s\t%3.2f?\t%s…\n",
					aircraft.icaoAddr, aircraft.callsign,
					sLatLon, isMlat, sAlt, metersInMiles(distance),
					durationSecondsElapsed(tPos))
			}
		}
	}
}
