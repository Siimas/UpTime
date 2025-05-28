package logger

import (
	"context"
	"encoding/json"
	"log"

	"uptime/internal/cache"
	"uptime/internal/constants"
	"uptime/internal/events"
	"uptime/internal/models"
	"uptime/internal/postgres"
	"uptime/internal/util/workers"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, pooldb *pgxpool.Pool, kc *events.KafkaConsumer, rdb *redis.Client) {
	defer log.Println("⚠️ - Logger Shutting Down")

	wp := workers.NewPool(constants.LoggerChanSize, false, func(msg *kafka.Message) {
		log.Printf("📋 Logging Monitor with key %s", msg.Key)
		handleMonitorResult(ctx, msg, pooldb, rdb)
	})

	wp.Launch(constants.LoggerWorkerCount, ctx)

	log.Println("✅ - Logger Online")

	kc.Subscribe(ctx, []string{constants.KafkaMonitorResultsTopic}, func(ev *kafka.Message) {
		wp.Dispach(ev)
	})
}

func handleMonitorResult(ctx context.Context, km *kafka.Message, pooldb *pgxpool.Pool, rdb *redis.Client) {
	var monitorResult models.MonitorResult
	if err := json.Unmarshal(km.Value, &monitorResult); err != nil {
		log.Printf("🚨 Error converting json to monitor result: %s\n", err)
		return
	}

	if err := cache.UpdateMonitorStatus(ctx, monitorResult, rdb); err != nil {
		log.Printf("🚨 Error updating monitor status: %s\n", err)
	}

	if err := postgres.StoreMonitorResult(ctx, monitorResult, pooldb); err != nil {
		log.Printf("🚨 Error storing monitor result (postgres): %s\n", err)
	}
}
