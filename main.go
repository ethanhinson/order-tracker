package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
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

// server is used to implement StaysServer.
type server struct {
	service.UnimplementedDeliveryTrackerServer
}

func (s *server) Track(ctx context.Context, input *service.TrackDelivery) (*service.DeliveryStatus, error) {
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

	return &service.DeliveryStatus{
		OnTime:       false,
		ExpectedTime: arrival.String(),
	}, err
}

var RedisConnection *redis.Client

func InitializeRedisConnection() {
	var conn = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	var _, err = conn.Ping().Result()

	if err != nil {
		log.Panicf("Could not connect to redis: %s", err)
	}

	RedisConnection = conn
}


func main() {
	readConfig()
	InitializeRedisConnection()
	defer RedisConnection.Close()
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