package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func SetMongoClient(c *mongo.Client) {
	client = c
}

func GetMongoClient() *mongo.Client {
	return client
}
