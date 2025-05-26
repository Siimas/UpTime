package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"uptime/internal/cache"
	"uptime/internal/constants"
	"uptime/internal/events"
	"uptime/internal/models"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, db *pgx.Conn, kc *events.KafkaConsumer, rdb *redis.Client) {
	defer log.Println("⚠️ - Scheduler Shutting Down")
	log.Println("✅ - Scheduler Online")

	kc.Subscribe(ctx, []string{constants.KafkaMonitorScheduleTopic}, func(ev *kafka.Message) {
		handleMonitorScheduler(ctx, ev, rdb)
	})

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
