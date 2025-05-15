package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO: IMPROVE THIS
func NewConsumer(bootstrapServers string) (*kafka.Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": bootstrapServers,
		// "sasl.username":     "<CLUSTER API KEY>",
		// "sasl.password":     "<CLUSTER API SECRET>",

		// Fixed properties
		// "security.protocol": "SASL_SSL",
		// "sasl.mechanisms":   "PLAIN",
		"group.id":          "kafka-go-local-consumer",
		"auto.offset.reset": "earliest"})
		if err != nil {
			return nil, err
		}
	return c, nil
}

func NewProducer(bootstrapServers string) (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": bootstrapServers,

		// Fixed properties
		"acks": "all",
	})

	if err != nil {
		return nil, err
	}

	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	return p, nil
}
