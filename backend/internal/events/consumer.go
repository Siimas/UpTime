package events

import (
	"context"
	"log"
	"time"
	"uptime/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	Consumer *kafka.Consumer
}

type HandlerFunc func(ev *kafka.Message)

func (kc *KafkaConsumer) Subscribe(ctx context.Context, topics []string, handler HandlerFunc) {
	defer kc.Consumer.Close()

	if err := kc.Consumer.SubscribeTopics(topics, nil); err != nil {
		log.Fatalf("⚠️ Error Subscribing to topic: %s\n", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("🔴 Consumer stopped")
			return
		default:
			ev, err := kc.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() != kafka.ErrTimedOut {
					log.Fatalf("⚠️ Logger Kafka error: %s\n", kafkaErr)
				}
				continue
			}

			handler(ev)
		}
	}
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
