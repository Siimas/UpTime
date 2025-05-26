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
	"uptime/internal/logger"
	"uptime/internal/scheduler"
	"uptime/internal/util/color"
	"uptime/internal/worker"

	"uptime/internal/cache"
	"uptime/internal/events"
	"uptime/internal/postgres"
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

	redisClient := cache.NewClient(config.RedisURL)
	defer redisClient.Close()

	loggerConsumer := events.NewLocalConsumer()
	defer loggerConsumer.Consumer.Close()

	schedulerConsumer := events.NewLocalConsumer()
	defer schedulerConsumer.Consumer.Close()

	kafkaProducer := events.NewLocalProducer()
	defer kafkaProducer.Producer.Close()

	if err := cache.SeedRedisFromPostgres(ctx, db, redisClient); err != nil {
		log.Println("🚨 Error starting flusher: " + err.Error())
	}

	go logger.Run(ctx, pooldb, loggerConsumer, redisClient)

	go worker.Run(ctx, redisClient, kafkaProducer)

	go scheduler.Run(ctx, db, schedulerConsumer, redisClient)

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
