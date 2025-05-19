package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uptime/internal/constants"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func RunMonitorResults(ctx context.Context, db *pgx.Conn, kc *kafka.Consumer, rdb *redis.Client) {
	err := kc.SubscribeTopics([]string{constants.KafkaMonitorResultsTopic}, nil)
	if err != nil {
		fmt.Printf("Couldn't subscribe to topic: %s\n", err)
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
			
			go handleMonitorResult(ctx, ev, db, rdb)
		}
	}

	kc.Close()
}

func handleMonitorResult(ctx context.Context, km *kafka.Message, db *pgx.Conn, rdb *redis.Client) error {
	var monitorResult MonitorResult
	if err := json.Unmarshal(km.Value, &monitorResult); err != nil {
		fmt.Printf("Error converting json to monitor result: %s\n", err)
		return err
	}

	if err := LogMonitorResult(monitorResult); err != nil {
		fmt.Printf("Error logging monitor result: %s\n", err)
		return err
	}

	if err := StoreMonitorResult(ctx, monitorResult, db); err != nil {
		fmt.Printf("Error storing monitor result: %s\n", err)
		return err
	}

	if monitorResult.Status == StatusOffline {
		if err := UpdateMonitorStatus(ctx, monitorResult, rdb); err != nil {
			fmt.Printf("Error updating monitor status: %s\n", err)
			return err
		}
	}

	return nil
}