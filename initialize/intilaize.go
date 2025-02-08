package initialize

import (
	"context"
	"fmt"
	"github.com/omniful/go_commons/log"

	"github.com/omniful/go_commons/config" // Replace with the actual config package you use
	oredis "github.com/omniful/go_commons/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"oms/pkg/redis"
	//"log"
	"time"
)

func InitializeLog(ctx context.Context) {
	err := log.InitializeLogger(
		log.Formatter(config.GetString(ctx, "log.format")),
		log.Level(config.GetString(ctx, "log.level")),
	)
	if err != nil {
		log.WithError(err).Panic("unable to initialise log")
	}

}

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

func InitializeRedis() {
	//r := oredis.NewClient(&oredis.Config{
	//	ClusterMode:  config.GetBool(ctx, "redis.clusterMode"),
	//	Hosts:        config.GetStringSlice(ctx, "redis.hosts"),
	//	DB:           config.GetUint(ctx, "redis.db"),
	//	DialTimeout:  10 * time.Second,
	//	ReadTimeout:  10 * time.Second,
	//	WriteTimeout: 10 * time.Second,
	//	IdleTimeout:  10 * time.Second,
	//})

	config := &oredis.Config{
		Hosts:       []string{"localhost:6379"},
		PoolSize:    50,
		MinIdleConn: 10,
	}

	client := oredis.NewClient(config)
	//defer client.Close()
	fmt.Println("Initialized Redis Client")
	redis.SetClient(client)
}
