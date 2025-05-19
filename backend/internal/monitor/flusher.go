package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"uptime/internal/constants"
	"uptime/internal/util"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func RunMonitorFlusher(ctx context.Context, db *pgx.Conn, kc *kafka.Consumer, rdb *redis.Client) {
	if err := LoadMonitorConfigs(ctx, db, rdb); err != nil {
		fmt.Println("Error starting flusher: " + err.Error())
	}

	err := kc.SubscribeTopics([]string{constants.KafkaMonitorActionTopic}, nil)
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

			go handleMonitorFlusher(ctx, ev, rdb)
		}
	}

	kc.Close()
}

func handleMonitorFlusher(ctx context.Context, km *kafka.Message, rdb *redis.Client) error {
	var monitorEvent MonitorEvent
	if err := json.Unmarshal(km.Value, &monitorEvent); err != nil {
		fmt.Printf("Error converting json to monitor action: %s\n", err)
		return err
	}

	util.PrettyPrint(monitorEvent)

	switch monitorEvent.Action {
	case MonitorDelete:
		if err := DeleteMonitor(ctx, monitorEvent.Monitor.Id, rdb); err != nil {
			fmt.Printf("Error deleting monitor (%s): %s\n", monitorEvent.Monitor.Id, err)
			return err
		}
	default:
		if err := ScheduleMonitor(ctx, monitorEvent.Monitor, rdb); err != nil {
			fmt.Printf("Error %s monitor (%s): %s\n", monitorEvent.Action.string(), monitorEvent.Monitor.Id, err)
			return err
		}
	}

	return nil
}

func LoadMonitorConfigs(ctx context.Context, db *pgx.Conn, rdb *redis.Client) error {
	monitors, err := GetActiveMonitors(ctx, db)
	if err != nil {
		fmt.Println("Failed to load monitors: " + err.Error())
	}

	var wg sync.WaitGroup

	for _, m := range monitors {
		wg.Add(1)

		go func(m Monitor) {
			defer wg.Done()
			if err := ScheduleMonitor(ctx, m, rdb); err != nil {
				fmt.Printf("Failed schedule monitor (%s): %s", m.Id, err.Error())
			}
			fmt.Println("Monitor Scheduled!")
			util.PrettyPrint(m)
		}(m)
	}

	wg.Wait()
	return nil
}
