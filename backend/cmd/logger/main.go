package main

import (
	"context"
	"log"
	"uptime/internal/cache"
	"uptime/internal/config"
	"uptime/internal/events"
	"uptime/internal/logger"
	"uptime/internal/postgres"
	"uptime/internal/util/color"
)

func main() {
	log.Println(color.Colorize(color.Blue, `
	 _   _     _____ _             _                           
	| | | |_ _|_   _(_)_ __  ___  | |   ___  __ _ __ _ ___ _ _ 
	| |_| | '_ \| | | | '  \/ -_) | |__/ _ \/ _  / _  / -_) '_|
	 \___/| .__/|_| |_|_|_|_\___| |____\___/\__, \__, \___|_|  
	      |_|                               |___/|___/         
	`))

	log.Println(color.Colorize(color.Green, "Logging service is starting..."))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.Load()

	pooldb := postgres.NewPoolConnection(ctx, config.PostgresURL)
	defer pooldb.Close()

	loggerConsumer := events.NewCloudConsumer()
	defer loggerConsumer.Consumer.Close()

	redisClient := cache.NewClient(config.RedisURL)
	defer redisClient.Close()

	logger.Run(ctx, pooldb, loggerConsumer, redisClient)
}
