package monitor

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"uptime/internal/constants"
	"uptime/internal/models"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func Ping(monitorId string, monitor models.MonitorCache, kp *kafka.Producer) {
	if monitor.Endpoint == "" {
		log.Println("🚨 Empty Endpoint! ")
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

	messageData, err := json.Marshal(monitorResult)
	if err != nil {
		log.Println("🚨 Error marshaling data:", err)
		return
	}

	topic := constants.KafkaMonitorResultsTopic
	key := constants.RedisMonitorKey + ":" + monitorId
	kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(messageData),
	}, nil)
}
