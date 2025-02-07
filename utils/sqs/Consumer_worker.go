package sqs

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/sqs"
	"log"
	psqs "oms/pkg/sqs"
)

type ExampleHandler struct{}

func (h *ExampleHandler) Process(ctx context.Context, message *[]sqs.Message) error {
	//TODO implement me
	//panic("implement me")

	for _, msg := range *message {
		err := h.Handle(&msg)
		if err != nil {

		}
	}
	return nil
}

func (h *ExampleHandler) Handle(msg *sqs.Message) error {
	fmt.Println("Processing message:", string(msg.Value))

	return nil
}

func StartConsumerWorker(ctx context.Context) {

	// Set up consumer
	handler := &ExampleHandler{}
	consumer, err := sqs.NewConsumer(
		psqs.QueueGlobal,
		1, // Number of workers
		1, // Concurrency per worker
		handler,
		10,    // Max messages count
		30,    // Visibility timeout
		false, // Is async
		false, // Send batch message
	)

	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}

	consumer.Start(ctx)

	// Let the consumer run for a while
	//time.Sleep(10 * time.Second)

	//consumer.Close()
}
