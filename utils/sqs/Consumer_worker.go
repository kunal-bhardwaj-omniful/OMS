package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	//"github.com/google/uuid"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/http"
	//"gorm.io/datatypes"

	//"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
	"github.com/omniful/go_commons/sqs"
	"log"

	"oms/domain/models"
	psqs "oms/pkg/sqs"
	"time"
)

const getBySkuId = "/sku/"

var (
// service ser1.Service
)

//var (
//	once                  sync.Once
//	wmsServiceInstance    *Client
//	wmsServiceInstanceErr error
//)

//	type Client struct {
//		*http.Client
//	}
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

		fmt.Println(records)

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

		for _, v := range records {

			sku_id := v[1]
			hub_id := v[2]

			fmt.Printf("sku_id %s and hub_id %s", sku_id, hub_id)

			res, err1 := GetSku(ctx, sku_id, client)
			//fmt.Println("FSDFKJFD")
			fmt.Println(err1)
			if err1 != nil {
				fmt.Println("2")

				fmt.Printf("Error from wms service %w", err)
				fmt.Println("id  not  present in wms service")
				//return
			}

			data, ok := res.([]byte)
			if !ok {
				fmt.Println("Error: res is not of type []byte")
				//return
			}

			// Define a map to store the decoded JSON data
			var result map[string]interface{}

			// Unmarshal JSON into the map
			err := json.Unmarshal(data, &result)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				//return
			}

			// Print the decoded JSON data
			fmt.Println("Decoded JSON:", result)

			fmt.Println("res from isvc ", res)
			//response.NewSuccessResponse(res)
			//fmt.Println("body", json.Unmarshal( []byte(res),&sfsd )

			// cyclic dependency
			//order1.SaveOrderinDb(models.Order{
			//	HubID: hub_id,
			//	SkuId: sku_id,
			//	Qty:   0,
			//})

		}

		//fmt.Println(m[f])
	}

	return nil
}

func GetSku(ctx context.Context, skuID string, client *interservice_client.Client) (res interface{}, err *interservice_client.Error) {
	//request := &http.Request{
	//	Url: fmt.Sprintf("/sku/%s", skuID),
	//}

	//var skuData interface{} // Use proper struct instead of `interface{}`
	//response, err := client.Get(request, &skuData)

	fmt.Println("ISVC START")

	res, err = client.Execute(ctx,
		http.APIGet,
		&http.Request{
			Url: fmt.Sprintf("/sku/%s", skuID),
			//Headers: headers,
		},
		&res,
	)

	fmt.Println("ISVC END")

	if err != nil {
		return
	}

	return
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
