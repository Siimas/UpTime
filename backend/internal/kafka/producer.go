package kafka

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewProducer(broker string) *kafka.Producer {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"acks":              "all",
	})

	if err != nil {
		log.Fatalf("Failed to create kafka consumer: %s", err)
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					// fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
					// 	*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	return p
}