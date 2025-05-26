package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime/internal/cache"
	"uptime/internal/constants"
	"uptime/internal/events"
	"uptime/internal/models"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, db *pgx.Conn, kc *events.KafkaConsumer, rdb *redis.Client) {
	log.Println("✅ - Scheduler Online")
	defer log.Println("⚠️ - Scheduler Shutting Down")
	ctx.Done()

	err := kc.Consumer.SubscribeTopics([]string{constants.KafkaMonitorScheduleTopic}, nil)
	if err != nil {
		log.Printf("🚨 Couldn't subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := kc.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() != kafka.ErrTimedOut {
					log.Printf("⚠️ Scheduler Kafka error: %s\n", kafkaErr)
					os.Exit(1)
				}
				continue
			}

			go handleMonitorScheduler(ctx, ev, rdb)
		}
	}

	kc.Consumer.Close()
}

func handleMonitorScheduler(ctx context.Context, km *kafka.Message, rdb *redis.Client) error {
	var monitorEvent models.MonitorEvent
	if err := json.Unmarshal(km.Value, &monitorEvent); err != nil {
		log.Printf("Error converting json to monitor action: %s\n", err)
		return err
	}

	switch monitorEvent.Action {
	case models.MonitorCreate:
		if err := cache.ScheduleMonitor(ctx, monitorEvent.Monitor, rdb); err != nil {
			log.Printf("Error %s monitor (%s): %s\n", monitorEvent.Action.String(), monitorEvent.Monitor.Id, err)
			return err
		}
	case models.MonitorDelete:
		if err := cache.DeleteMonitor(ctx, monitorEvent.Monitor.Id, rdb); err != nil {
			log.Printf("Error deleting monitor (%s): %s\n", monitorEvent.Monitor.Id, err)
			return err
		}
	default:

	}

	return nil
}
