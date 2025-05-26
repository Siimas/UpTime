package logger

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
	"uptime/internal/postgres"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, pooldb *pgxpool.Pool, kc *events.KafkaConsumer, rdb *redis.Client) {
	log.Println("✅ - Logger Online")
	defer log.Println("⚠️ - Logger Shutting Down")

	if err := kc.Consumer.SubscribeTopics([]string{constants.KafkaMonitorResultsTopic}, nil); err != nil {
		log.Fatalf("Couldn't subscribe to topic: %s\n", err)
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
					log.Fatalf("⚠️ Logger Kafka error: %s\n", kafkaErr)
				}
				continue
			}

			go handleMonitorResult(ctx, ev, pooldb, rdb)
		}
	}

	kc.Consumer.Close()
}

func handleMonitorResult(ctx context.Context, km *kafka.Message, pooldb *pgxpool.Pool, rdb *redis.Client) {
	var monitorResult models.MonitorResult
	if err := json.Unmarshal(km.Value, &monitorResult); err != nil {
		log.Printf("🚨 Error converting json to monitor result: %s\n", err)
		return
	}

	log.Println("📋 Logging Monitor:", monitorResult)

	if err := postgres.StoreMonitorResult(ctx, monitorResult, pooldb); err != nil {
		log.Printf("🚨 Error storing monitor result (postgres): %s\n", err)
	}

	if err := cache.UpdateMonitorStatus(ctx, monitorResult, rdb); err != nil {
		log.Printf("🚨 Error updating monitor status: %s\n", err)
	}

}
