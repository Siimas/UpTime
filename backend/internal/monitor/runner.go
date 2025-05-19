package monitor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"uptime/internal/constants"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
)

func RunMonitorRunner(ctx context.Context, rdb *redis.Client, kp *kafka.Producer) {
	for {
		now := float64(time.Now().Unix())

		monitorIDs, err := rdb.ZRangeByScore(ctx, constants.RedisMonitorsScheduleKey, &redis.ZRangeBy{
			Min:   "-inf",
			Max:   fmt.Sprintf("%f", now),
			Count: 1000,
		}).Result()
		if err != nil {
			log.Printf("failed to fetch due monitors: %v", err)
			continue
		}

		for _, key := range monitorIDs {
			monitor, err := GetMonitor(ctx, rdb, key)
			if err != nil {
				log.Printf("error retrieving monitor %s: %v", key, err)
				continue
			}

			go Ping(strings.Split(key, ":")[1], monitor, kp)

			nextPing := time.Now().Add(time.Duration(monitor.Interval) * time.Second).Unix()
			rdb.ZAdd(ctx, constants.RedisMonitorsScheduleKey, redis.Z{
				Score:  float64(nextPing),
				Member: key,
			})
		}
	}
}
