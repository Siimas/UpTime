package monitor

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"uptime/internal/constants"

	"github.com/jackc/pgx/v5"
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

// todo
func LogMonitorResult(mr MonitorResult) error {
	fmt.Println("Monitor Result: {")
	fmt.Printf("\tId: %s\n", mr.Id)
	fmt.Printf("\tStatus: %s\n", mr.Status)
	fmt.Printf("\tLatency: %d ms\n", mr.Latency)
	fmt.Printf("\tDate: %s\n", mr.Date)
	fmt.Println("}")
	return nil
}

func UpdateMonitorStatus(ctx context.Context, mr MonitorResult, rdb *redis.Client) error {
	key := constants.RedisMonitorKey + ":" + mr.Id
	return rdb.HSet(ctx, key, "status", mr.Status).Err()
}

func StoreMonitorResult(ctx context.Context, mr MonitorResult, db *pgx.Conn) error {
	sql := `INSERT INTO monitor_results (monitor_id, status, latency_ms, checked_at) VALUES ($1, $2, $3, $4)`
	
	if _, err := db.Exec(ctx, sql, mr.Id, mr.Status, mr.Latency, mr.Date); err != nil {
		return err
	}

	return nil
}