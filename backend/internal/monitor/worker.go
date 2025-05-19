package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"uptime/internal/constants"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Ping(monitorId string, monitor MonitorCache, kp *kafka.Producer) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	var errorMessage string

	start := time.Now()
	resp, err := client.Head(monitor.Endpoint)
	latency := time.Since(start)
	if err != nil {
		fmt.Printf("Error pinging url: %s \n%s", monitor.Endpoint, err) // todo: improve
		errorMessage = err.Error()
	}
	defer func(e error) {
		if e != nil {
			resp.Body.Close()
		}
	}(err)

	fmt.Printf("Ping successful: %s (%d)\n", monitor.Endpoint, resp.StatusCode)

	var status MonitorStatus
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		status = StatusUp
	} else {
		status = StatusDown
	}

	monitorResult := MonitorResult{
		Id:      monitorId,
		Date:    time.Now().Format("2006-01-02 15:04:05-07"),
		Latency: latency.Milliseconds(),
		Status:  status,
		Code:    resp.StatusCode,
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
