package main

import (
	"context"
	"log"
	"uptime/internal/cache"
	"uptime/internal/config"
	"uptime/internal/events"
	"uptime/internal/util/color"
	"uptime/internal/worker"
)

func main() {
	log.Println(color.Colorize(color.Blue, `
	 _   _     _____ _            __      __       _           
	| | | |_ _|_   _(_)_ __  ___  \ \    / /__ _ _| |_____ _ _ 
	| |_| | '_ \| | | | '  \/ -_)  \ \/\/ / _ \ '_| / / -_) '_|
	 \___/| .__/|_| |_|_|_|_\___|   \_/\_/\___/_| |_\_\___|_|  
	      |_|                                                  
	`))

	log.Println(color.Colorize(color.Green, "Worker service is starting..."))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.Load()

	rdb := cache.NewClient(config.RedisURL)

	loggerProducer := events.NewCloudProducer()

	worker.Run(ctx, rdb, loggerProducer)
}
