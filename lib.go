package main

import (
	"fmt"
	. "github.com/waitr/tracker/service"
	"math"
	"os"
	"strconv"
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
 * formula, or ECEF based calculations. It is probably accurate enough when considering the use
 * case and context.
 */
func HaversineDistance(from *Point, to *Point) (distance float64) {
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
	t := time.Now().Add(time.Duration(((distance / speed) * 60 * 60) * 1000000000)).UTC()
	return t
}

const DateLayout = "2006-01-02T15:04:05Z"

func MakeTime(value string) time.Time {
	t, _ := time.Parse(DateLayout, value)
	return t
}


func HandleTrackingRequest(input *TrackDelivery, messenger SMSMessenger) (*DeliveryStatus, error) {
	speed, err := strconv.Atoi(os.Getenv("DRIVER_SPEED"))
	rate, err := strconv.Atoi(os.Getenv("DRIVER_RATE"))

	arrival := ArrivalTime(
		time.Now(),
		HaversineDistance(input.GetLocation(), input.GetDestination()),
		float64(speed),
		Unit(rate),
	)

	expected := MakeTime(input.GetArrivalTime())
	if arrival.Unix() <= expected.Unix() {
		return &DeliveryStatus{
			OnTime:       true,
			ExpectedTime: arrival.String(),
		}, err
	}

	// Respond async. In a real world context you would likely
	// queue these messages and send them "off main thread".
	// This is a common approach to reduce network latency
	// for the API and permit things like "retries".
	go (func() {
		sent, err := messenger.Send(SMSMessage{
			To:   input.GetContact(),
			From: os.Getenv("TWILIO_PHONE"),
			Body: fmt.Sprintf("Order with id: %s will be late. Expected time: %s", input.GetOrderId(), arrival.String()),
		})
		if !sent || err != nil {
			fmt.Printf("Message for late order (%s) was not sent with err: %v", input.GetOrderId(), err)
		}
	})()

	return &DeliveryStatus{
		OnTime:       false,
		ExpectedTime: arrival.String(),
	}, err
}