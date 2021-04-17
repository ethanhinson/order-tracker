package main

import (
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

// server is used to implement StaysServer.
type server struct {
	service.UnimplementedDeliveryTrackerServer
}


func main() {
	readConfig()
	InitializeRedisConnection()
	defer RedisConnection.Close()
	lis, err := net.Listen("tcp", os.Getenv("TRACKER_SERVICE_PORT"))
	log.Printf("Preparing to listen for connections...")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Creating new server instance...")
	var s = grpc.NewServer()
	log.Printf("Registering server instance...")
	service.RegisterDeliveryTrackerServer(s, &server{})
	log.Printf("Registered server instance...start tracker service...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}