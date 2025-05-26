package events

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// todo: work on this
func SubscribeTopic(
	ctx context.Context,
	kc *kafka.Consumer,
	topics []string,
	handler func(ctx context.Context, ev *kafka.Message),
) {
	err := kc.SubscribeTopics(topics, nil)
	if err != nil {
		fmt.Printf("Couldn't subscribe to topic: %s\n", err)
		os.Exit(1)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := kc.ReadMessage(100 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() != kafka.ErrTimedOut {
					fmt.Printf("⚠️ Kafka error: %s\n", kafkaErr)
				}
				continue
			}
			go handler(ctx, ev)
		}
	}

	kc.Close()
}
