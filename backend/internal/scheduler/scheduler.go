package scheduler

import (
	"context"
	"encoding/json"
	"log"
	"uptime/internal/cache"
	"uptime/internal/constants"
	"uptime/internal/events"
	"uptime/internal/models"
	"uptime/internal/postgres"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, db *pgx.Conn, kc *events.KafkaConsumer, rdb *redis.Client) {
	defer log.Println("⚠️ - Scheduler Shutting Down")
	log.Println("✅ - Scheduler Online")

	kc.Subscribe(ctx, []string{constants.KafkaMonitorScheduleTopic}, func(ev *kafka.Message) {
		handleMonitorScheduler(ctx, ev, db, rdb)
	})

}

func handleMonitorScheduler(ctx context.Context, km *kafka.Message, db *pgx.Conn, rdb *redis.Client) error {
	var monitorEvent models.MonitorEvent
	if err := json.Unmarshal(km.Value, &monitorEvent); err != nil {
		log.Printf("Error converting json to monitor action: %s\n", err)
		return err
	}

	log.Printf("🔄 Scheduling Monitor: [%s]\t%s \n", monitorEvent.Action, monitorEvent.MonitorId)

	switch monitorEvent.Action {
	case models.MonitorCreate:
		monitor, err := postgres.GetSingleMonitor(ctx, db, monitorEvent.MonitorId)
		if err != nil {
			log.Printf("Error getting monitor (%s): %s\n", monitorEvent.MonitorId, err)
			return err
		}

		if err := cache.ScheduleMonitor(ctx, monitor, rdb); err != nil {
			log.Printf("Error schedulling monitor (%s): %s\n", monitorEvent.MonitorId, err)
			return err
		}

	case models.MonitorUpdate:
		monitor, err := postgres.GetSingleMonitor(ctx, db, monitorEvent.MonitorId)
		if err != nil {
			log.Printf("Error getting monitor (%s): %s\n", monitorEvent.MonitorId, err)
			return err
		}

		if monitor.Active {
			if err := cache.ScheduleMonitor(ctx, monitor, rdb); err != nil {
				log.Printf("Error schedulling monitor (%s): %s\n", monitorEvent.MonitorId, err)
				return err
			}
		} else {
			if err := cache.DeleteMonitor(ctx, monitorEvent.MonitorId, rdb); err != nil {
				log.Printf("Error deleting monitor (%s): %s\n", monitorEvent.MonitorId, err)
				return err
			}
		}

	case models.MonitorDelete:
		if err := cache.DeleteMonitor(ctx, monitorEvent.MonitorId, rdb); err != nil {
			log.Printf("Error deleting monitor (%s): %s\n", monitorEvent.MonitorId, err)
			return err
		}
	default:

	}

	log.Printf("✅ Successfully Schedulled Monitor: [%s]\t%s \n", monitorEvent.Action, monitorEvent.MonitorId)
	return nil
}
