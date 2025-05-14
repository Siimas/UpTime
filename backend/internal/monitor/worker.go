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
	}
	defer resp.Body.Close()


	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		fmt.Printf("Ping successful: %s (%d)\n", monitor.Endpoint, resp.StatusCode)

		topic := constants.KafkaMonitorResultsTopic
		key := "monitor:"+monitor.Id
		monitorResult := MonitorResult{
			Id: monitor.Id,
			Date: time.Now().String(),
			Latency: latency.Milliseconds(),
			Status: monitor.Status,
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
}
