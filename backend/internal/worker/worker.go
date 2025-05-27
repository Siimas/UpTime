package worker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"uptime/internal/cache"
	"uptime/internal/constants"
	"uptime/internal/events"
	"uptime/internal/models"

	"github.com/redis/go-redis/v9"
)

type scheduleTask struct {
	Key string
	Wg  *sync.WaitGroup
}

func Run(ctx context.Context, rdb *redis.Client, kp *events.KafkaProducer) {
	log.Println("✅ - Monitor Runner Online")
	defer log.Println("⚠️ - Monitor Runner Shutting Down")

	pingChan := make(chan string, 1000)
	scheduleChan := make(chan scheduleTask, 1000)

	pingWorkerCount := 10
	for i := range pingWorkerCount {
		go pingWorker(i, ctx, pingChan, rdb, kp)
	}

	shceduleWorkerCount := 10
	for i := range shceduleWorkerCount {
		go scheduleWorker(i, ctx, scheduleChan, rdb)
	}

	var wg sync.WaitGroup

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

		wg.Add(len(monitorIDs))

		for _, key := range monitorIDs {
			pingChan <- key
			scheduleChan <- scheduleTask{Key: key, Wg: &wg}
		}

		wg.Wait()
	}
}

func pingWorker(
	id int,
	ctx context.Context,
	pingChan <-chan string,
	rdb *redis.Client,
	kp *events.KafkaProducer,
) {
	log.Printf("👷 Worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("👷 Worker %d shutting down", id)
			return
		case key := <-pingChan:
			monitor, err := cache.GetMonitor(ctx, rdb, key)
			if err != nil {
				log.Printf("Error retrieving monitor %s: %v", key, err)
				continue
			}

			monitorId := strings.Split(key, ":")[1]

			Ping(monitorId, monitor, kp)
		}
	}
}

func scheduleWorker(id int, ctx context.Context, scheduleChan <-chan scheduleTask, rdb *redis.Client) {
	log.Printf("👷 Schedule Worker %d started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("👷 Schedule Worker %d shutting down", id)
			return
		case task := <-scheduleChan:
			monitor, err := cache.GetMonitor(ctx, rdb, task.Key)
			if err != nil {
				log.Printf("Error retrieving monitor %s: %v", task.Key, err)
				continue
			}

			nextPing := time.Now().Add(time.Duration(monitor.Interval) * time.Second).Unix()
			rdb.ZAdd(ctx, constants.RedisMonitorsScheduleKey, redis.Z{
				Score:  float64(nextPing),
				Member: task.Key,
			})

			task.Wg.Done()
		}
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
