package main

import (
	"context"
	"fmt"
	"os"
	"uptime/internal/kafka"
	"uptime/internal/monitor"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	kc, err := kafka.NewConsumer("localhost:" + os.Getenv("KAFKA_PLAINTEXT_PORTS"))
	if err != nil {
		fmt.Printf("Failed to create kafka producer: %s", err)
		os.Exit(1)
	}

	monitor.RunMonitorResults(ctx, kc)
}
