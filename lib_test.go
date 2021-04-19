package main

import (
	. "github.com/waitr/tracker/service"
	"math"
	"testing"
	"time"
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
	distance := HaversineDistance(glenwoodSprings, aspen)
	if math.Round(distance) != 43 {
		t.Fail()
	}
	// I think real tests/fuzz tests would choose points starting
	// from near the poles and measure how much distances start
	// to vary
}

func TestRadians(t *testing.T) {
	r := Radians(0)
	if r != 0 {
		t.Fail()
	}
	r1 := Radians(180)
	if math.Round(r1) != math.Round(math.Pi) {
		t.Fail()
	}
}

func TestArrivalTime(t *testing.T) {
	now := time.Now()
	at := ArrivalTime(now, 1, 1, KMH)
	if at.Hour() != now.Hour() + 1 {
		t.Fail()
	}
	atMPH := ArrivalTime(now, 1, 1, MPH)
	if atMPH.Hour() != now.Hour() + 1 {
		t.Fail()
	}
}

func TestMakeTime(t *testing.T) {

}