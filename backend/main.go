package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"
	"uptime/internal"

	"github.com/redis/go-redis/v9"
)

type Monitor struct {
	Endpoint string `json:"endpoint"`
	Interval int    `json:"interval"`
	Status   string `json:"status"`
}

type MonitorStatus int

func (ms MonitorStatus) String() string {
	return [...]string{"Online", "Offline"}[ms]
}

var ctx = context.Background()

func main() {

	rdb := internal.BuildRedisClient()

	for {
		now := float64(time.Now().Unix())

		monitorIDs, err := rdb.ZRangeByScore(ctx, "monitors_schedule", &redis.ZRangeBy{
			Min:   "-inf",                 // From the beginning of time
			Max:   fmt.Sprintf("%f", now), // Until now
			Count: 1000,
		}).Result()
		if err != nil {
			log.Printf("failed to fetch due monitors: %v", err)
		}

		for _, key := range monitorIDs {
			log.Printf("Found monitor with key: %s", key)

			monitorJson, err := rdb.HGetAll(ctx, key).Result()
			if err == redis.Nil {
				fmt.Println("Key does not exist")
				break
			} else if err != nil {
				fmt.Println("Error:", err)
				break
			}

			var monitor Monitor
			monitor.Endpoint = monitorJson["endpoint"]
			monitor.Interval, _ = strconv.Atoi(monitorJson["interval"])
			monitor.Status = monitorJson["status"]

			go func(monitor Monitor) {
				fmt.Printf("Monitor: %+v\n", monitor)

			}(monitor)

			nextPing := time.Now().Add(time.Duration(monitor.Interval) * time.Second).Unix()
			rdb.ZAdd(ctx, "monitors_schedule", redis.Z{
				Score:  float64(nextPing),
				Member: key,
			}).Result()
			// if err != nil {
			// 	fmt.Println("Failed to reschedule:", err)
			// }
		}

	}

}
