package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	//"github.com/google/uuid"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/http"
	//"gorm.io/datatypes"

	//"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
	"github.com/omniful/go_commons/sqs"
	"log"
	"oms/domain/models"
	pmongo "oms/pkg/mongo"
	psqs "oms/pkg/sqs"
	"oms/utils/kafka"
	"time"
)

const getBySkuId = "/sku/"

//var (
//	once                  sync.Once
//	wmsServiceInstance    *Client
//	wmsServiceInstanceErr error
//)

//	type Client struct {
//		*http.Client
//	}

type SKUResponse struct {
	Data    SKUData `json:"data"`
	Message string  `json:"message"`
	Status  string  `json:"status"`
}

// SKUData represents the data object with only the ID field.
type SKUData struct {
	ID string `json:"id"`
}

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
	//fmt.Println("Processing message:", string(msg.Value))
	fmt.Println("Processing message:", string(msg.Value))

	var b models.BulkOrderEvent

	// Unmarshal the message body into the struct
	err := json.Unmarshal(msg.Value, &b)
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

		//headers := []string{"file path"}
		//m := records.ToMaps(headers)
		// Process the records

		//fmt.Println(records)

		ctx := context.TODO()
		//client, err := NewClient(ctx)

		if err != nil {
			log.Println("error in inter service client", err)
		}

		config := interservice_client.Config{
			ServiceName: "oms-service",
			BaseURL:     "http://localhost:8081/api/v1",
			Timeout:     5 * time.Second,
		}

		client, err := interservice_client.NewClientWithConfig(config)
		if err != nil {
			panic(err)
		}

		// Todo write invalid orders in new csv
		//
		//
		//dest := csv.Destination{}
		//dest.SetFileName("rejectedOrder.csv")
		//wr, _ := csv.NewCommonCSVWriter()
		//wr.SetDestination(dest)

		for _, v := range records {

			sku_id := v[0]
			hub_id := v[1]

			fmt.Printf("sku_id %s and hub_id %s", sku_id, hub_id)

			res, err1 := GetSku(ctx, sku_id, client)
			//fmt.Println("FSDFKJFD")
			//fmt.Println(err1)
			if err1 == nil {
				// VALID
				fmt.Println("process SKU valid ", res.Data)
			} else {
				// invalid sku

				continue
			}

			res, err1 = GetHub(ctx, hub_id, client)
			//fmt.Println("FSDFKJFD")
			//fmt.Println(err1)
			if err1 == nil {
				// VALID
				fmt.Println("process hub valid valid ", res.Data)
			} else {
				// invalid hub
				continue
			}

			//reached this line means valid hubid and skuid

			// put in mongoDb with status onHold

			fmt.Println("VALID ORDER")

			//save order to dB
			client := pmongo.GetMongoClient()

			order := &models.Order{
				HubID:  hub_id,
				SkuId:  sku_id,
				Qty:    0,
				Status: "onHold",
			}
			err = SaveOrder(client, order)

			if err != nil {
				fmt.Println("ERROR IN SAVING MDB ", err)
			}
			//push order event to kafka
			kafka.PushOrderToKafka(order)

		}

		//fmt.Println(m[f])
	}

	return nil
}

func SaveOrder(client *mongo.Client, order *models.Order) error {
	collection := client.Database("orders").Collection("orders")

	// Ensure default status is set
	if order.Status == "" {
		order.Status = "onhold"
	}

	// Insert into MongoDB
	res, err := collection.InsertOne(context.Background(), order)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	fmt.Println("saved in mongo and id ", res.InsertedID)
	return nil
}

func GetHub(ctx context.Context, hubID string, client *interservice_client.Client) (*SKUResponse, *interservice_client.Error) {

	//fmt.Println("ISVC START")

	var skuRes SKUResponse

	res, err := client.Get(
		&http.Request{
			Url: fmt.Sprintf("/hub/%s", hubID),
			//Headers: headers,
		}, &skuRes)

	if err != nil {
		return &skuRes, err
	}

	//fmt.Println(res.StatusCode())

	err1 := json.Unmarshal(res.Body(), &skuRes)
	if err1 != nil {
		return &skuRes, &interservice_client.Error{
			Message:    "",
			Errors:     nil,
			StatusCode: 404,
			Data:       nil,
		}
	}

	fmt.Println(skuRes)

	fmt.Println("ISVC END")

	if err != nil {
		return &skuRes, err
	}

	return &skuRes, nil
}

func GetSku(ctx context.Context, skuID string, client *interservice_client.Client) (*SKUResponse, *interservice_client.Error) {

	//fmt.Println("ISVC START")

	var skuRes SKUResponse

	res, err := client.Get(
		&http.Request{
			Url: fmt.Sprintf("/sku/%s", skuID),
			//Headers: headers,
		}, &skuRes)

	if err != nil {
		return &skuRes, err
	}

	//fmt.Println(res.StatusCode())

	err1 := json.Unmarshal(res.Body(), &skuRes)
	if err1 != nil {
		return &skuRes, &interservice_client.Error{
			Message:    "",
			Errors:     nil,
			StatusCode: 404,
			Data:       nil,
		}
	}

	fmt.Println(skuRes)

	fmt.Println("ISVC END")

	if err != nil {
		return &skuRes, err
	}

	return &skuRes, nil
}

//func NewClient(ctx context.Context) (*Client, error) {
//	once.Do(func() {
//		baseUrl := config.GetString(ctx, "wms.baseUrl")
//		timeout := config.GetInt64(ctx, "wms.timeout")
//
//		httpClient, err := http.NewHTTPClient("wms-service", baseUrl, &http2.Transport{
//			MaxIdleConns:       10,
//			IdleConnTimeout:    60 * time.Second,
//			DisableCompression: true,
//		}, http.WithTimeout(time.Duration(timeout)*time.Second))
//		if err != nil {
//			wmsServiceInstanceErr = err
//			return
//		}
//
//		wmsServiceInstance = &Client{
//			httpClient,
//		}
//	})
//
//	return wmsServiceInstance, wmsServiceInstanceErr
//}

//func (sku *Client) GetSkuByID(ctx context.Context, sku_id string) (err *interservice_client.Error) {
//	requestEndpoint := getBySkuId + sku_id
//	fmt.Println(requestEndpoint)
//	//var response models.MapsAPIResponse
//	//headers := utils.GetDefaultHeaders(ctx)
//
//	//queryParams["placeid"] = []string{placeID}
//	//queryParams["key"] = []string{config.GetString(ctx, "googleApi.key")}
//	var data interface{}
//	result, clientErr := sku.Client.Get(&http.Request{Url: requestEndpoint}, data)
//	if clientErr != nil {
//		err = &interservice_client.Error{Message: fmt.Sprintf("error while calling sku :: err :: %v", clientErr.Error())}
//		log.Println(result.Error())
//		return
//	}
//
//	//
//	//result.Status()
//	//
//	//clientErr = json.Unmarshal(result.Body(), &response)
//	//if clientErr != nil {
//	//	err = &interservice_client.Error{Message: fmt.Sprintf("error while calling sku :: err :: %v", clientErr.Error())}
//	//	log.Println(result.Error())
//	//	return
//	//}
//
//	if result.Status() == "200" {
//		return nil
//	} else {
//		err = &interservice_client.Error{Message: "GetSKU Error : Invalid Request"}
//	}
//	return
//}

func StartConsumerWorker(ctx context.Context) {
	fmt.Println("sqs consumer created")
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

}
