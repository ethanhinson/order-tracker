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
)

// Ensure that we can read the environment variables.
func readConfig() {
	_ = gotenv.Load()
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