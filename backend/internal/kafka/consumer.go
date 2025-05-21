package kafka

import (
	"log"
	"uptime/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func NewConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
		"sasl.username":      config.GetEnv("KAFKA_SASL_USERNAME"),
		"sasl.password":      config.GetEnv("KAFKA_SASL_PASSWORD"),
		"security.protocol":  "SASL_SSL",
		"sasl.mechanisms":    "PLAIN",
		"group.id":           config.GetEnv("KAFKA_GROUP_ID"),
		"client.id":          config.GetEnv("KAFKA_CLIENT_ID"),
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": true,
	})
	if err != nil {
		log.Fatalf("Failed to create kafka producer: %s", err)
	}
	return c
}
