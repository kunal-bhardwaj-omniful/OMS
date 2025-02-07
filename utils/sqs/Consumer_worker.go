package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/sqs"
	"log"
	"oms/domain/models"
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

	var b models.BulkOrderEvent

	// Unmarshal the message body into the struct
	err := json.Unmarshal([]byte(msg.Value), &b)
	if err != nil {
		fmt.Println("Error unmarshaling message:", err)
		return err
	}

	csvReader, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(b.FilePath),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)

	if err != nil {
		log.Fatal(err)
	}
	err = csvReader.InitializeReader(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	for !csvReader.IsEOF() {
		var records csv.Records
		records, err = csvReader.ReadNextBatch()
		if err != nil {
			log.Fatal(err)
		}
		//records.ToMaps()
		// Process the records
		fmt.Println(records)
	}

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
