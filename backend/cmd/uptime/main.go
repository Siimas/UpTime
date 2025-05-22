package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime/internal/config"
	"uptime/internal/http"
	"uptime/internal/util/color"

	"uptime/internal/kafka"
	"uptime/internal/monitor"
	"uptime/internal/postgres"
	"uptime/internal/redisclient"
)

func main() {
	log.Println(color.Colorize(color.Blue, `
  _   _     _____ _           
 | | | |_ _|_   _(_)_ __  ___ 
 | |_| | '_ \| | | | '  \/ -_)
  \___/| .__/|_| |_|_|_|_\___|
       |_|                    
	`))
	log.Println("Uptime service is starting...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown handler
	go handleShutdown(cancel)

	config := config.Load()

	db := postgres.NewConnection(ctx)
	defer db.Close(context.Background())

	pooldb := postgres.NewPoolConnection(ctx, config.PostgresURL)
	defer pooldb.Close()

	redisClient := redisclient.NewClient(config.RedisURL)
	defer redisClient.Close()

	kafkaConsumer := kafka.NewConsumer()
	defer kafkaConsumer.Close()

	kafkaProducer := kafka.NewProducer()
	defer kafkaProducer.Close()

	if err := redisclient.SeedRedisFromPostgres(ctx, db, redisClient); err != nil {
		log.Println("🚨 Error starting flusher: " + err.Error())
	}

	go monitor.RunMonitorFlusher(ctx, db, kafkaConsumer, redisClient)

	go monitor.RunMonitorResults(ctx, pooldb, kafkaConsumer, redisClient)

	go monitor.RunMonitorRunner(ctx, redisClient, kafkaProducer)

	go http.StartServer(ctx, config.HTTPServerAddr)

	// Block main until context is cancelled
	<-ctx.Done()
	log.Println("Uptime service shutting down.")
}

// handleShutdown cancels the context on SIGINT/SIGTERM
func handleShutdown(cancel context.CancelFunc) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	log.Println("Shutdown signal received.")
	cancel()

	// Give time for cleanup
	time.Sleep(1 * time.Second)
}
