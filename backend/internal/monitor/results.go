package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"uptime/internal/constants"
	"uptime/internal/models"
	"uptime/internal/redisclient"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func RunMonitorResults(ctx context.Context, pooldb *pgxpool.Pool, kc *kafka.Consumer, rdb *redis.Client) {
	log.Println("✅ - Monitor Results Online")
	defer log.Println("⚠️ - Monitor Results Shutting Down")

	if err := kc.SubscribeTopics([]string{constants.KafkaMonitorResultsTopic}, nil); err != nil {
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
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := kc.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() != kafka.ErrTimedOut {
					fmt.Printf("Kafka error: %s\n", kafkaErr)
				}
				continue
			}

			go handleMonitorResult(ctx, ev, pooldb, rdb)
		}
	}

	kc.Close()
}

func handleMonitorResult(ctx context.Context, km *kafka.Message, pooldb *pgxpool.Pool, rdb *redis.Client) {
	var monitorResult models.MonitorResult
	if err := json.Unmarshal(km.Value, &monitorResult); err != nil {
		fmt.Printf("🚨 Error converting json to monitor result: %s\n", err)
		return
	}

	if err := StoreMonitorResult(ctx, monitorResult, pooldb); err != nil {
		fmt.Printf("🚨 Error storing monitor result (postgres): %s\n", err)
	}

	if err := redisclient.UpdateMonitorStatus(ctx, monitorResult, rdb); err != nil {
		fmt.Printf("🚨 Error updating monitor status: %s\n", err)
	}

}
