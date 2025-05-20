package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO: IMPROVE THIS
func NewConsumer(broker string, groupId string) *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		// "sasl.username":     "<CLUSTER API KEY>",
		// "sasl.password":     "<CLUSTER API SECRET>",
		// "security.protocol": "SASL_SSL",
		// "sasl.mechanisms":   "PLAIN",
		"group.id":          groupId,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create kafka producer: %s", err)
	}
	return c
}
