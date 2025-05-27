package main

import (
	"context"
	"log"
	"uptime/internal/config"
	"uptime/internal/events"
	"uptime/internal/http"
	"uptime/internal/postgres"
	"uptime/internal/util/color"
)

func main() {
	log.Println(color.Colorize(color.Blue, `
	 _   _     _____ _               _   ___ ___ 
	| | | |_ _|_   _(_)_ __  ___    /_\ | _ \_ _|
	| |_| | '_ \| | | | '  \/ -_)  / _ \|  _/| | 
	 \___/| .__/|_| |_|_|_|_\___| /_/ \_\_| |___|
	      |_|                                    
	`))

	log.Println(color.Colorize(color.Green, "API service is starting..."))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.Load()

	db := postgres.NewConnection(ctx)
	defer db.Close(context.Background())

	kafkaProducer := events.NewLocalProducer()
	defer kafkaProducer.Producer.Close()

	server := http.NewServer(config.HTTPServerAddr, db, kafkaProducer)
	server.Run(ctx)
}
