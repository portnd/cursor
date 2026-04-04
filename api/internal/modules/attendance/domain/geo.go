package domain

import "math"

const earthRadiusMeters = 6371000.0

// HaversineDistanceMeters returns the great-circle distance between two WGS84 points in meters.
func HaversineDistanceMeters(lat1, lon1, lat2, lon2 float64) float64 {
	if lat1 == lat2 && lon1 == lon2 {
		return 0
	}
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	Δφ := (lat2 - lat1) * math.Pi / 180
	Δλ := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMeters * c
}
