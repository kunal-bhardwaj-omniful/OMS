package main

import (
	"context"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/shutdown"
	"go.mongodb.org/mongo-driver/mongo"
	"oms/initialize"
	psqs "oms/pkg/sqs"
	"oms/router"
	"oms/utils/kafka"
	"oms/utils/sqs"
	"os"
	"time"
)

var client *mongo.Client

func main() {
	// Initialize config
	os.Setenv("CONFIG_SOURCE", "local")

	//fmt.Println(os.Getwd())
	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initialising config, err: %v", err)
		panic(err)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Error while getting context from config, err: %v", err)
		panic(err)
	}

	//context = context.Background()
	client, err = initialize.ConnectMongo(ctx)
	if err != nil {
		log.Panicf("Error while running mongo client, err: %v", err)
	}

	psqs.IntiializeSqs(ctx)
	sqs.StartConsumerWorker(ctx)

	kafka.InitializeKafka()
	go kafka.StartConsumerKafka()
	func() {
		time.Sleep(time.Millisecond * 100)
		kafka.PushOrderToKafka()

	}()
	//time.Sleep(time.Second * 2)
	runHttp(ctx)

}

func runHttp(ctx context.Context) {

	server := http.InitializeServer(config.GetString(ctx, "server.port"), 10*time.Second, 10*time.Second, 70*time.Second)

	//err := router.Initialize(ctx, server)
	//if err != nil {
	//	log.Errorf(err.Error())
	//	panic(err)
	//}

	err := router.InternalRoutes(ctx, server, client)
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	log.Infof("Starting " + config.GetString(ctx, "name") + " server on port" + config.GetString(ctx, "server.port"))

	err = server.StartServer(config.GetString(ctx, "name"))
	if err != nil {
		log.Errorf(err.Error())
		panic(err)
	}

	<-shutdown.GetWaitChannel()
}
