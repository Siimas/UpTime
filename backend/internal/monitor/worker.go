package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"uptime/internal/constants"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Ping(monitor Monitor, kp *kafka.Producer) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()
	resp, err := client.Head(monitor.Endpoint)
	latency := time.Since(start)
	if err != nil {
		fmt.Printf("Error pinging url: %s \n%s", monitor.Endpoint, err) // todo: improve
		return
	} 
	defer resp.Body.Close()

	fmt.Printf("Ping successful: %s (%d)\n", monitor.Endpoint, resp.StatusCode)

	topic := constants.KafkaMonitorResultsTopic
	key := "monitor:" + monitor.Id
	var status MonitorStatus
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		status = StatusOnline
	} else {
		status = StatusOffline
	}

	monitorResult := MonitorResult{
		Id:      monitor.Id,
		Date:    time.Now().Format("2006-01-02 15:04:05-07"),
		Latency: latency.Milliseconds(),
		Status:  status,
	}

	messageData, err := json.Marshal(monitorResult)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}

	kp.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(messageData),
	}, nil)
}
