package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
	"github.com/omniful/go_commons/kafka"
	"github.com/omniful/go_commons/pubsub"
	"oms/domain/models"
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
	fmt.Println("Received Kafka Message:", string(msg.Value))

	// Convert message byte slice to Order struct
	var order models.Order
	err := json.Unmarshal(msg.Value, &order)
	if err != nil {
		fmt.Println("Error decoding message:", err)
		return err
	}

	fmt.Println("Recovered Order Struct:", order)

	// Interservice client configuration
	config := interservice_client.Config{
		ServiceName: "oms-service",
		BaseURL:     "http://localhost:8081/api/v1",
		Timeout:     5 * time.Second,
	}

	client, err := interservice_client.NewClientWithConfig(config)
	if err != nil {
		fmt.Println("Error creating interservice client:", err)
		return err
	}

	// Fetch inventory details from WMS
	inventoryAvailable, err := CheckInventory(ctx, order.SkuId, order.HubID, order.Qty, client)
	if err != nil {
		fmt.Println("Error fetching inventory details:", err)
		return err
	}
	fmt.Println("Inventory fetched with val", inventoryAvailable)

	if inventoryAvailable {
		// Call WMS to decrease inventory
		err = DecreaseInventory(ctx, order.SkuId, order.HubID, order.Qty, client)
		if err != nil {
			fmt.Println("Error decreasing inventory:", err)
			return err
		}

		fmt.Println("Inventory successfully decreased for SKU:", order.SkuId)
	} else {
		fmt.Println("Insufficient inventory for SKU:", order.SkuId)
	}

	return nil
}
func CheckInventory(ctx context.Context, skuID, hubID string, qty int, client *interservice_client.Client) (bool, error) {
	var response struct {
		Data struct {
			AvailableQty int `json:"available_qty"`
		} `json:"data"`
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	res, err := client.Get(
		&http.Request{
			Url: fmt.Sprintf("/inventory?sku_id=%s&hub_id=%s", skuID, hubID),
		}, &response)

	if err != nil {
		return false, errors.New("inventory fetch error")
	}

	err1 := json.Unmarshal(res.Body(), &response)
	if err1 != nil {
		fmt.Println("marshal error")
	}

	fmt.Println()

	fmt.Println(response.Data.AvailableQty)

	if response.Data.AvailableQty >= qty {
		return true, nil
	}

	return false, nil
}

func DecreaseInventory(ctx context.Context, skuID, hubID string, qty int, client *interservice_client.Client) error {
	var response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	reqBody := map[string]interface{}{
		"sku_id":        skuID,
		"hub_id":        hubID,
		"available_qty": qty,
	}
	_, err := client.Post(
		&http.Request{
			Url:  "/inventory",
			Body: reqBody,
		}, &response)

	if err != nil {
		return errors.New("inventory update error")
	}

	if response.Status != "success" {
		return fmt.Errorf("failed to decrease inventory: %s", response.Message)
	}

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
