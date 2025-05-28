package cache

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"
	"uptime/internal/constants"
	"uptime/internal/models"
	"uptime/internal/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func GetMonitor(ctx context.Context, rdb *redis.Client, key string) (models.MonitorCache, error) {
	data, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return models.MonitorCache{}, err
	}
	if len(data) == 0 {
		return models.MonitorCache{}, redis.Nil
	}

	interval, err := strconv.Atoi(data["Interval"])
	if err != nil {
		interval = 60
	}

	status := models.MonitorStatus(data["Status"])

	return models.MonitorCache{
		Endpoint: data["Endpoint"],
		Interval: interval,
		Status:   status,
	}, nil
}

func UpdateMonitorStatus(ctx context.Context, mr models.MonitorResult, rdb *redis.Client) error {
	key := constants.RedisMonitorKey + ":" + mr.Id
	return rdb.HSet(ctx, key, "Status", mr.Status.String()).Err()
}

func DeleteMonitor(ctx context.Context, monitorId string, rdb *redis.Client) error {
	key := constants.RedisMonitorKey + ":" + monitorId

	if err := rdb.Del(ctx, key).Err(); err != nil {
		return err
	}

	if err := rdb.ZRem(ctx, constants.RedisMonitorsScheduleKey, key).Err(); err != nil {
		return err
	}

	return nil
}

func SeedRedisFromPostgres(ctx context.Context, db *pgx.Conn, rdb *redis.Client) error {
	log.Println("🌱 - Seeding Cache with Monitors")

	monitors, err := postgres.GetActiveMonitors(ctx, db)
	if err != nil {
		log.Println("🚨 Failed to load monitors: " + err.Error())
	}

	var wg sync.WaitGroup

	for _, m := range monitors {
		wg.Add(1)

		go func(m models.Monitor) {
			defer wg.Done()
			if err := ScheduleMonitor(ctx, m, rdb); err != nil {
				log.Printf("🚨 Failed schedule monitor (%s): %s", m.Id, err.Error())
			}
		}(m)
	}

	wg.Wait()
	return nil
}

func ScheduleMonitor(ctx context.Context, mr models.Monitor, rdb *redis.Client) error {
	// todo: check if not valid
	key := constants.RedisMonitorKey + ":" + mr.Id

	fields := map[string]any{
		"Endpoint": mr.Endpoint,
		"Status":   models.StatusDown.String(),
		"Interval": mr.Interval,
	}
	
	result, err := rdb.Exists(ctx, key).Result()
	switch {
	case err != nil:
		return err
	case result > 0:
		return nil
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
