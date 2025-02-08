package repo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"oms/domain/models"

	//"github.com/omniful/go_commons/db/sql/postgres"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

type Repository interface {
	SaveOrder(order models.Order) error
}

// Inter
type repository struct {
	db *mongo.Client
}

// Singleton instance of Repository and the sync.Once to ensure it's initialized only once.
var repo *repository
var repoOnce sync.Once

// NewRepository is the constructor function that ensures the Repository is initialized only once.
func NewRepository(db *mongo.Client) Repository {
	repoOnce.Do(func() {
		// Initialize the Repository with a given DbCluster.
		repo = &repository{
			db: db,
		}
	})
	return repo
}

func (r *repository) SaveOrder(order models.Order) error {
	// Access the database and collection
	collection := r.db.Database("orders").Collection("orders")

	// Create a document from the CSV record
	document := bson.D{
		//{"column1", record[0]},
		//{"column2", record[1]},
		//{"column3", record[2]},
		// Add more fields as needed based on your CSV structure
	}

	// Insert the document into MongoDB
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return fmt.Errorf("failed to insert document: %v", err)
	}

	return nil
}
