package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"uptime/internal/cache"
	"uptime/internal/constants"
	"uptime/internal/events"
	"uptime/internal/models"
	"uptime/internal/util/workers"

	"github.com/redis/go-redis/v9"
)

func Run(ctx context.Context, rdb *redis.Client, kp *events.KafkaProducer) {
	log.Println("✅ - Monitor Runner Online")
	defer log.Println("⚠️ - Monitor Runner Shutting Down")

	scheduleWp := workers.NewPool(constants.PingChanSize, true, func(key string) {
		monitor, err := cache.GetMonitor(ctx, rdb, key)
		if err != nil {
			log.Printf("Error retrieving monitor %s: %v", key, err)
		}

		nextPing := time.Now().Add(time.Duration(monitor.Interval) * time.Second).Unix()
		rdb.ZAdd(ctx, constants.RedisMonitorsScheduleKey, redis.Z{
			Score:  float64(nextPing),
			Member: key,
		})
	})

	scheduleWp.Launch(constants.PingWorkerCount, ctx)

	pingWp := workers.NewPool(constants.PingChanSize, false, func(key string) {
		monitor, err := cache.GetMonitor(ctx, rdb, key)
		if err != nil {
			log.Printf("Error retrieving monitor %s: %v", key, err)
		}

		monitorId := strings.Split(key, ":")[1]
		Ping(monitorId, monitor, kp)
	})

	pingWp.Launch(constants.PingWorkerCount, ctx)

	for {
		now := float64(time.Now().Unix())

		monitorIDs, err := rdb.ZRangeByScore(ctx, constants.RedisMonitorsScheduleKey, &redis.ZRangeBy{
			Min:   "-inf",
			Max:   fmt.Sprintf("%f", now),
			Count: 1000,
		}).Result()
		if err != nil {
			log.Printf("Failed to fetch due monitors: %v", err)
			continue
		}

		for _, key := range monitorIDs {
			pingWp.Dispach(key)
			scheduleWp.Dispach(key)
		}

		scheduleWp.Wg.Wait()
	}
}

func Ping(monitorId string, monitor models.MonitorCache, kp *events.KafkaProducer) {
	if monitor.Endpoint == "" {
		log.Println("🚨 Caught Empty Endpoint!")
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var statusCode int
	var errorMessage string
	var status models.MonitorStatus

	log.Printf("📡 Pinging: %s", monitor.Endpoint)

	start := time.Now()
	resp, err := client.Head(monitor.Endpoint)
	latency := time.Since(start)

	if err != nil {
		log.Printf("%s --> 🔴 Error: %s\n", monitor.Endpoint, err)
		errorMessage = err.Error()
	} else {
		defer resp.Body.Close()

		statusCode = resp.StatusCode

		log.Printf("%s --> 🟢 Ping successful (%d)", monitor.Endpoint, statusCode)

		if statusCode >= 200 && statusCode < 400 {
			status = models.StatusUp
		} else {
			status = models.StatusDown
		}
	}

	monitorResult := models.MonitorResult{
		Id:      monitorId,
		Date:    time.Now().Format("2006-01-02 15:04:05-07"),
		Latency: latency.Milliseconds(),
		Status:  status,
		Code:    statusCode,
		Error:   errorMessage,
	}

	topic := constants.KafkaMonitorResultsTopic
	key := constants.RedisMonitorKey + ":" + monitorId
	kp.ProduceMessage(topic, key, monitorResult)
}
