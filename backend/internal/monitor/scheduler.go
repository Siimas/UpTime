package monitor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func RunScheduler(ctx context.Context, rdb *redis.Client) {
	for {
		now := float64(time.Now().Unix())

		monitorIDs, err := rdb.ZRangeByScore(ctx, "monitors_schedule", &redis.ZRangeBy{
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

			go func(m Monitor) {
				fmt.Printf("Monitor: %+v\n", m)
			}(monitor)

			nextPing := time.Now().Add(time.Duration(monitor.Interval) * time.Second).Unix()
			rdb.ZAdd(ctx, "monitors_schedule", redis.Z{
				Score:  float64(nextPing),
				Member: key,
			})
		}
	}
}
