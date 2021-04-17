package main

import (
	. "github.com/waitr/tracker/service"
	"math"
	"time"
)

const earthRadius = 6371

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
 * between 2 points. I suspect that while this formula is the least accurate of the Vincenty
 * formula, and ECEF based calculations, that it is accurate enough when considering the use
 * case and context.
 */
func HaversineDistance(from Point, to Point) (distance float64) {
	return 2 * earthRadius * math.Asin(math.Sqrt(Haversine(Radians(to.GetLatitude()) - Radians(from.GetLatitude())) +
		math.Cos(Radians(from.GetLongitude())) * math.Cos(Radians(to.GetLongitude())) * Haversine(Radians(to.GetLongitude()) - Radians(from.GetLongitude()))))
}