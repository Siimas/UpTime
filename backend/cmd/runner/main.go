package main

import (
	"context"
	"fmt"
	"os"
	"uptime/internal/kafka"
	"uptime/internal/monitor"
	"uptime/internal/redis"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	kp, err := kafka.NewProducer("localhost:" + os.Getenv("KAFKA_PLAINTEXT_PORTS"))
	if err != nil {
		fmt.Printf("Failed to create kafka producer: %s", err)
		os.Exit(1)
	}

	rdb, err := redis.NewClient(os.Getenv("REDIS_URL"))
	if err != nil {
		fmt.Println("Couldn't establish connection with redis")
		os.Exit(1)
	}

	monitor.RunMonitorRunner(ctx, rdb, kp)
}
