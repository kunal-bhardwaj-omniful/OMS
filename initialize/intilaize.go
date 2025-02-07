package initialize

import (
	"context"
	"fmt"
	//"log"
	"time"

	"github.com/omniful/go_commons/config" // Replace with the actual config package you use
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo(ctx context.Context) (*mongo.Client, error) {
	// Read MongoDB configuration from the config package
	uri := config.GetString(ctx, "mongodb.uri")
	database := config.GetString(ctx, "mongodb.database")
	timeout := config.GetInt(ctx, "mongodb.timeout")

	if uri == "" || database == "" {
		return nil, fmt.Errorf("missing MongoDB config values")
	}

	// Set connection timeout
	clientOptions := options.Client().ApplyURI(uri)
	connCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(connCtx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	err = client.Ping(connCtx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("✅ Connected to MongoDB successfully!")
	return client, nil
}
