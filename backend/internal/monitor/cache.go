package monitor

import (
	"context"
	"strconv"
	"time"
	"uptime/internal/constants"

	"github.com/redis/go-redis/v9"
)

func GetMonitor(ctx context.Context, rdb *redis.Client, key string) (MonitorCache, error) {
	data, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return MonitorCache{}, err
	}
	if len(data) == 0 {
		return MonitorCache{}, redis.Nil
	}

	interval, err := strconv.Atoi(data["interval"])
	if err != nil {
		interval = 60
	}

	status := MonitorStatus(data["status"])

	return MonitorCache{
		Endpoint: data["endpoint"],
		Interval: interval,
		Status:   status,
	}, nil
}

func UpdateMonitorStatus(ctx context.Context, mr MonitorResult, rdb *redis.Client) error {
	key := constants.RedisMonitorKey + ":" + mr.Id
	return rdb.HSet(ctx, key, "status", mr.Status.string()).Err()
}

func ScheduleMonitor(ctx context.Context, mr Monitor, rdb *redis.Client) error {
	key := constants.RedisMonitorKey + ":" + mr.Id

	fields := map[string]interface{}{
		"Endpoint": mr.Endpoint,
		"Status":   StatusDown.string(),
		"Interval": mr.Interval,
	}
	if err := rdb.HSet(ctx, key, fields).Err(); err != nil {
		return err
	}

	nextPing := time.Now().Unix()
	if err := rdb.ZAdd(ctx, constants.RedisMonitorsScheduleKey, redis.Z{
		Score:  float64(nextPing),
		Member: key,
	}).Err(); err != nil {
		return err
	}

	return nil
}
