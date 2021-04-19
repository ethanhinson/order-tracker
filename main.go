package main

import (
	"context"
	"fmt"
	"github.com/subosito/gotenv"
	"github.com/waitr/tracker/service"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// Ensure that we can read the environment variables.
func readConfig() {
	_ = gotenv.Load()
}

func HandleTrackingRequest(input *service.TrackDelivery, messenger SMSMessenger) (*service.DeliveryStatus, error) {
	speed, err := strconv.Atoi(os.Getenv("DRIVER_SPEED"))
	rate, err := strconv.Atoi(os.Getenv("DRIVER_RATE"))

	arrival := ArrivalTime(
		time.Now(),
		HaversineDistance(input.GetLocation(), input.GetDestination()),
		float64(speed),
		Unit(rate),
	)

	if arrival.Unix() <= input.GetArrivalTime().GetSeconds() {
		return &service.DeliveryStatus{
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

	return &service.DeliveryStatus{
		OnTime:       false,
		ExpectedTime: arrival.String(),
	}, err
}

type server struct {
	service.UnimplementedDeliveryTrackerServer
}

func (s *server) Track(ctx context.Context, input *service.TrackDelivery) (*service.DeliveryStatus, error) {
	msg := new(TwilioMessenger)
	resp, err := HandleTrackingRequest(input, msg)
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
	return resp, err
}

func main() {
	readConfig()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("TRACKER_SERVICE_PORT")))
	log.Printf("Preparing to listen for connections...")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var s = grpc.NewServer()
	log.Printf("Registering server instance...")
	service.RegisterDeliveryTrackerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}