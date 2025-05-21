package monitor

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime/internal/constants"
	"uptime/internal/models"
	"uptime/internal/redisclient"
	"uptime/internal/util"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func RunMonitorFlusher(ctx context.Context, db *pgx.Conn, kc *kafka.Consumer, rdb *redis.Client) {
	log.Println("✅ - Monitor Flusher Online")
	defer log.Println("⚠️ - Monitor Flusher Shutting Down")

	err := kc.SubscribeTopics([]string{constants.KafkaMonitorActionTopic}, nil)
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
			ev, err := kc.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() != kafka.ErrTimedOut {
					log.Printf("-Kafka error: %s\n", kafkaErr)
				}
				continue
			}

			go handleMonitorFlusher(ctx, ev, rdb)
		}
	}

	kc.Close()
}

func handleMonitorFlusher(ctx context.Context, km *kafka.Message, rdb *redis.Client) error {
	var monitorEvent models.MonitorEvent
	if err := json.Unmarshal(km.Value, &monitorEvent); err != nil {
		log.Printf("Error converting json to monitor action: %s\n", err)
		return err
	}

	util.PrettyPrint(monitorEvent)

	switch monitorEvent.Action {
	case models.MonitorDelete:
		if err := redisclient.DeleteMonitor(ctx, monitorEvent.Monitor.Id, rdb); err != nil {
			log.Printf("Error deleting monitor (%s): %s\n", monitorEvent.Monitor.Id, err)
			return err
		}
	default:
		if err := redisclient.ScheduleMonitor(ctx, monitorEvent.Monitor, rdb); err != nil {
			log.Printf("Error %s monitor (%s): %s\n", monitorEvent.Action.String(), monitorEvent.Monitor.Id, err)
			return err
		}
	}

	return nil
}
