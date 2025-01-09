package main

import (
	"math"
)

func cprNLFunction(lat float64) byte {
	lat = abs(lat) // Ensure latitude is positive for comparison

	latThresholds := []float64{
		10.47047130, 14.82817437, 18.18626357, 21.02939493, 23.54504487, 25.82924707,
		27.93898710, 29.91135686, 31.77209708, 33.53993436, 35.22899598, 36.85025108,
		38.41241892, 39.92256684, 41.38651832, 42.80914012, 44.19454951, 45.54626723,
		46.86733252, 48.16039128, 49.42776439, 50.67150166, 51.89342469, 53.09516153,
		54.27817472, 55.44378444, 56.59318756, 57.72747354, 58.84763776, 59.95459277,
		61.04917774, 62.13216659, 63.20427479, 64.26616523, 65.31845310, 66.36171008,
		67.39646774, 68.42322022, 69.44242631, 70.45451075, 71.45986473, 72.45884545,
		73.45177442, 74.43893416, 75.42056257, 76.39684391, 77.36789461, 78.33374083,
		79.29428225, 80.24923213, 81.19801349, 82.13956981, 83.07199445, 83.99173563,
		84.89166191, 85.75541621, 86.53536998, 87.00000000,
	}

	// Iterate over the thresholds to find the appropriate NL value
	for i, threshold := range latThresholds {
		if lat < threshold {
			return byte(59 - i)
		}
	}

	// Default case if latitude exceeds all thresholds
	return 1
}

// Helper function to calculate the absolute value of a float
func abs(val float64) float64 {
	if val < 0 {
		return -val
	}
	return val
}

func cprNFunction(lat float64, fflag bool) byte {
	// Descriptive variable for offset based on fflag
	fflagOffset := byte(1)
	if !fflag {
		fflagOffset = 0
	}

	// Calculate the adjusted NL value
	return calculateNL(lat, fflagOffset)
}

// Helper function to calculate NL value, encapsulating logic for NL adjustments
func calculateNL(lat float64, offset byte) byte {
	nl := cprNLFunction(lat) - offset
	if nl < 1 {
		nl = 1
	}
	return nl
}

func cprDlonFunction(lat float64, fflag bool, surface bool) float64 {
	var sfc float64
	if surface {
		sfc = 90.0
	} else {
		sfc = 360.0
	}

	return sfc / float64(cprNFunction(lat, fflag))

}

func decodeAC12Field(ac12Data uint) int32 {
	q := (ac12Data & 0x10) == 0x10
	if q {
		n := int32((ac12Data&0x0FE0)>>1) + int32(ac12Data&0x000F)
		return (n * 25) - 1000
	}
	/* TODO
	// Make N a 13 bit Gillham coded altitude by inserting M=0 at bit 6
	int n = ((AC12Field & 0x0FC0) << 1) |
					 (AC12Field & 0x003F);
	n = ModeAToModeC(decodeID13Field(n));
	if (n < -12) {
			return INVALID_ALTITUDE;
	}
	return (100 * n);
	*/
	return int32(math.MaxInt32)
}

const EarthRadiusMeters = 6371e3 // Earth's radius in meters

// degToRad converts degrees to radians.
func degToRad(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

// GreatCircle calculates the great-circle distance between two geographic points in meters.
func GreatCircle(lat1Deg, lon1Deg, lat2Deg, lon2Deg float64) float64 {
	lat1Rad := degToRad(lat1Deg)
	lon1Rad := degToRad(lon1Deg)
	lat2Rad := degToRad(lat2Deg)
	lon2Rad := degToRad(lon2Deg)

	// Handle edge case where points are very close to avoid NaN due to floating-point precision issues
	if math.Abs(lat1Rad-lat2Rad) < 0.0001 && math.Abs(lon1Rad-lon2Rad) < 0.0001 {
		return 0.0
	}

	// Calculate great-circle distance using the haversine formula
	return EarthRadiusMeters * math.Acos(
		math.Sin(lat1Rad)*math.Sin(lat2Rad)+
			math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Cos(math.Abs(lon1Rad-lon2Rad)),
	)
}

func metersInMiles(dist float64) float64 {
	return dist / float64(1609.34721869)
}
