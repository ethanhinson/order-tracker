package main

import (
	. "github.com/waitr/tracker/service"
	"math"
	"time"
)

// We are using KM
const (
	earthRadius = 6371
	km = 1.60934
)

// We'll support miles or kilometers when calculating arrival
// time.
type Unit int

const (
	KMH = iota
	MPH
)

// Convert degrees to radians
func Radians(deg float64) float64 {
	return deg * math.Pi / 180
}

// "Half versed sine"
func Haversine(t float64) float64 {
	return .5 * (1 - math.Cos(t))
}

/**
 * I researched the various approaches to calculating the geodesic distance
 * between 2 points. I suspect that while this formula is the less accurate when compared to the Vincenty
 * formula, or ECEF based calculations that. It is probably accurate enough when considering the use
 * case and context.
 */
func HaversineDistance(from Point, to Point) (distance float64) {
	return 2 * earthRadius * math.Asin(math.Sqrt(Haversine(Radians(to.GetLatitude()) - Radians(from.GetLatitude())) +
		math.Cos(Radians(from.GetLongitude())) * math.Cos(Radians(to.GetLongitude())) * Haversine(Radians(to.GetLongitude()) - Radians(from.GetLongitude()))))
}

// Convert MPH -> KMH
func mphTokmh(mph float64) float64 {
	return mph * km
}

// Given a distance in KM, a rate of speed, and a `unit`. Calculate the
// time in which the distance would be traveled.
func ArrivalTime(start time.Time, distance float64, speed float64, unit Unit) time.Time {
	if start.IsZero() {
		start = time.Now()
	}
	if unit == MPH {
		speed = mphTokmh(speed)
	}
	t := time.Now().Add(time.Duration(((distance / speed) * 60 * 60) * 1000000000))
	return t
}