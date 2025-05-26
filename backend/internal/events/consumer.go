package events

import (
	"context"
	"log"
	"uptime/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

// todo: add possiblity to add handler functions for topics
func (kc *KafkaConsumer) Consume(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("🔴 Consumer stopped")
				return
			default:
				ev := kc.Consumer.Poll(100)
				switch e := ev.(type) {
				case *kafka.Message:
					log.Printf("⬇️ Received message: %s = %s\n", string(e.Key), string(e.Value))
				case kafka.Error:
					log.Printf("⚠️ Consumer Kafka error: %v\n", e)
				}
			}
		}
	}()
}

func NewCloudConsumer() *KafkaConsumer {
	kc := createConsumer(kafka.ConfigMap{
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

	return &KafkaConsumer{kc}
}

func NewLocalConsumer() *KafkaConsumer {
	kc := createConsumer(kafka.ConfigMap{
		"bootstrap.servers": config.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
		"group.id":          config.GetEnv("KAFKA_GROUP_ID"),
		"auto.offset.reset": "earliest",
	})

	return &KafkaConsumer{kc}
}

func createConsumer(config kafka.ConfigMap) *kafka.Consumer {
	c, err := kafka.NewConsumer(&config)
	if err != nil {
		log.Fatalf("Failed to create kafka consumer: %s", err)
	}

	return c
}
