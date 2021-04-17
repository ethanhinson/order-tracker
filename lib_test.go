package main

import (
	. "github.com/waitr/tracker/service"
	"testing"
)

func TestHaversineDistance(t *testing.T) {
	glenwoodSprings := &Point{
		Latitude: 39.5505,
		Longitude: 107.3248,
	}
	aspen := &Point{
		Latitude: 39.1911,
		Longitude: 106.8175,
	}
	distance := HaversineDistance(*glenwoodSprings, *aspen)
	if distance != 43.25771880484925 {
		t.Fail()
	}
}

func TestRadians(t *testing.T) {
	r := Radians(0)
	if r != 0 {
		t.Fail()
	}
}