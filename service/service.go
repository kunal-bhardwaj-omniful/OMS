package service

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/log"
	"oms/domain/models"
	"oms/repo"
	"oms/utils/sqs"
	"os"
	"time"
)

type Service interface {
	ProcessOrder(filePath string) error
	SaveOrderinDb(order models.Order) error
}

type service struct {
	repo repo.Repository
}

// NewService is the constructor function to create a new instance of ConcreteService.
func NewService(r repo.Repository) Service {
	return &service{
		repo: r,
	}
}

func (s *service) ProcessOrder(filePath string) error {
	// Check if the file exists
	_, err := os.Stat(filePath)

	fmt.Println(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist")
		}
		return fmt.Errorf("unable to access file: %v", err)
	}

	bulkOrderEvent := models.BulkOrderEvent{
		FilePath: filePath,
		User: models.User{
			ID:        "1",
			FirstName: "Kunal",
			LastName:  "Bhardwaj",
			Email:     "kb@gmail.con",
		},
		RequestTime: time.Now(),
	}

	err = sqs.PushEmailMessageToSQS(context.Background(), bulkOrderEvent)
	if err != nil {
		log.Errorf("sqs push error %w", err)
	}
	// Open the CSV file
	//file, err := os.Open(filePath)
	//if err != nil {
	//	return fmt.Errorf("unable to open file: %v", err)
	//}
	//defer file.Close()
	//

	// Read and parse the CSV content
	//reader := csv.NewReader(file)
	//records, err := reader.ReadAll()
	//if err != nil {
	//	return fmt.Errorf("unable to read CSV content: %v", err)
	//}

	// Process each record and save it to MongoDB
	//for _, record := range records {
	//	// Save the record to MongoDB via the repository
	//	err := s.repo.SaveOrder(record)
	//	if err != nil {
	//		return fmt.Errorf("unable to save record to MongoDB: %v", err)
	//	}
	//}

	return nil
}

func (s *service) SaveOrderinDb(order models.Order) error {
	err := s.repo.SaveOrder(order)
	return err
}
