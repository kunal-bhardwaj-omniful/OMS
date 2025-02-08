package kafka

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
	"time"
)

// Implement message handler
type MessageHandler struct{}

func (h *MessageHandler) Process(ctx context.Context, message *pubsub.Message) error {
	//TODO implement me
	fmt.Println("processing kafka msg")
	err := h.Handle(ctx, message)
	if err != nil {

	}

	return nil

	//panic("implement me")
}

func (h *MessageHandler) Handle(ctx context.Context, msg *pubsub.Message) error {
	// Process message
	fmt.Println("rec kafka msg ", msg.Value)
	return nil
}

func StartConsumerKafka() {
	// Initialize consumer with configuration
	consumer := kafka.NewConsumer(
		kafka.WithBrokers([]string{"localhost:9092"}),
		kafka.WithConsumerGroup("my-consumer-group"),
		kafka.WithClientID("my-consumer"),
		kafka.WithKafkaVersion("2.8.1"),
		kafka.WithRetryInterval(time.Second),
		//kafka.WithDeadLetterConfig(&kafka.DeadLetterQueueConfig{
		//	Queue:     "dlq-queue",
		//	Account:   "aws-account",
		//	Region:    "us-east-1",
		//	ShouldLog: true,
		//}),
	)
	defer consumer.Close()

	// Register message handler for topic
	handler := &MessageHandler{}
	consumer.RegisterHandler("my-topic", handler)

	// Start consuming messages
	ctx := context.Background()
	fmt.Println("consumer kafka started")
	consumer.Subscribe(ctx)
}
