package main

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"log"
	"os"
	"time"
)

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