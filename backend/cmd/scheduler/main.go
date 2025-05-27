package main

import (
	"context"
	"log"
	"uptime/internal/cache"
	"uptime/internal/config"
	"uptime/internal/events"
	"uptime/internal/postgres"
	"uptime/internal/scheduler"
	"uptime/internal/util/color"
)

func main() {
	log.Println(color.Colorize(color.Blue, `
	 _   _     _____ _             ___     _           _      _         
	| | | |_ _|_   _(_)_ __  ___  / __| __| |_  ___ __| |_  _| |___ _ _ 
	| |_| | '_ \| | | | '  \/ -_) \__ \/ _| ' \/ -_) _\ | || | / -_) '_|
	 \___/| .__/|_| |_|_|_|_\___| |___/\__|_||_\___\__,_|\_,_|_\___|_|  
	      |_|                                                           
	`))

	log.Println(color.Colorize(color.Green, "Scheduler service is starting..."))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.Load()

	db := postgres.NewConnection(ctx)
	defer db.Close(context.Background())

	schedulerConsumer := events.NewLocalConsumer()
	defer schedulerConsumer.Consumer.Close()

	redisClient := cache.NewClient(config.RedisURL)
	defer redisClient.Close()

	if err := cache.SeedRedisFromPostgres(ctx, db, redisClient); err != nil {
		log.Println("🚨 Error starting flusher: " + err.Error())
	}

	scheduler.Run(ctx, db, schedulerConsumer, redisClient)
}
