package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"uptime/internal/constants"
	"uptime/internal/models"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Ping(monitorId string, monitor models.MonitorCache, kp *kafka.Producer) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var statusCode int
	var errorMessage string
	var status models.MonitorStatus

	log.Println("Pinging: ", monitor)

	start := time.Now()
	resp, err := client.Head(monitor.Endpoint)
	latency := time.Since(start)
	if err != nil {
		fmt.Printf("Error pinging url: %s \n%s", monitor.Endpoint, err)
		errorMessage = err.Error()
	} else {
		defer resp.Body.Close()

		statusCode = resp.StatusCode

		fmt.Printf("Ping successful: %s (%d)\n", monitor.Endpoint, statusCode)

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
		fmt.Println("Error marshaling data:", err)
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
