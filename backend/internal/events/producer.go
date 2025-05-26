package events

import (
	"log"
	"uptime/internal/config"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
}

func (kp *KafkaProducer) ProduceMessage(topic, key, value string) error {
	deliveryChan := make(chan kafka.Event)

	err := kp.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          []byte(value),
	}, deliveryChan)

	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		log.Fatalf("delivery failed: %v", m.TopicPartition.Error)
	}

	log.Println("✅ delivered message to", m.TopicPartition)
	return nil
}

func NewCloudProducer() *KafkaProducer {
	kp := createProducer(kafka.ConfigMap{
		"bootstrap.servers": config.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
		"sasl.username":     config.GetEnv("KAFKA_SASL_USERNAME"),
		"sasl.password":     config.GetEnv("KAFKA_SASL_PASSWORD"),
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"client.id":         config.GetEnv("KAFKA_CLIENT_ID"),
		"acks":              "all",
	})

	return &KafkaProducer{kp}
}

func NewLocalProducer() *KafkaProducer {
	kp := createProducer(kafka.ConfigMap{
		"bootstrap.servers": config.GetEnv("KAFKA_BOOTSTRAP_SERVERS"),
	})
	return &KafkaProducer{kp}
}

func createProducer(config kafka.ConfigMap) *kafka.Producer {
	p, err := kafka.NewProducer(&config)
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
