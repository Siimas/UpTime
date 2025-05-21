package kafka

import (
	"log"
	"uptime/internal/config"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO: IMPROVE THIS
func NewConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.GetEnv("KAFKA_BOOTSRAP_SERVERS"),
		"sasl.username":     config.GetEnv("KAFKA_SASL_USERNAME"),
		"sasl.password":     config.GetEnv("KAFKA_SASL_PASSWORD"),
		"security.protocol": config.GetEnv("KAFKA_SECURITY_PROTOCOL"),
		"sasl.mechanisms":   config.GetEnv("KAFKA_SASL_MECHANISM"),
		"group.id":          config.GetEnv("KAFKA_GROUP_ID"),
	})
	if err != nil {
		log.Fatalf("Failed to create kafka producer: %s", err)
	}
	return c
}
