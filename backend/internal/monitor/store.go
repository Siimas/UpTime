package monitor

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

func GetMonitor(ctx context.Context, rdb *redis.Client, key string) (Monitor, error) {
	data, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return Monitor{}, err
	}
	if len(data) == 0 {
		return Monitor{}, redis.Nil
	}

	interval, err := strconv.Atoi(data["interval"])
	if err != nil {
		interval = 60
	}

	status := MonitorStatus(data["status"])

	return Monitor{
		Id: strings.Split(key, ":")[1],
		Endpoint: data["endpoint"],
		Interval: interval,
		Status:   status,
	}, nil
}

func LogMonitorResult(mr MonitorResult) error {
	fmt.Println("Monitor Result: {")
	fmt.Printf("\tId: %s\n", mr.Id)
	fmt.Printf("\tStatus: %s\n", mr.Status)
	fmt.Printf("\tLatency: %d ms\n", mr.Latency)
	fmt.Printf("\tDate: %s\n", mr.Date)
	fmt.Println("}")
	return nil
}

func StoreMonitorResult(mr MonitorResult) error {
	return nil
}