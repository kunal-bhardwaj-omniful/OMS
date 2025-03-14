package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/omniful/go_commons/config"

	//"github.com/omniful/go_commons/config"
	error2 "github.com/omniful/go_commons/error"
	"github.com/omniful/go_commons/log"
	"github.com/omniful/go_commons/sqs"
	log1 "log"
	"oms/domain/models"
	error3 "oms/pkg/error"
	psqs "oms/pkg/sqs"
)

func PushEmailMessageToSQS(
	ctx context.Context,
	req models.BulkOrderEvent,
) (cusErr error2.CustomError) {

	_, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
	})
	if err != nil {
		log1.Fatalf("Failed to create AWS session: %v", err)
	}

	log.Infof("sqs sess created ")
	msg, err := json.Marshal(req)
	if err != nil {

		cusErr = error2.NewCustomError(error3.MarshalError, err.Error())
		return
	}
	//
	//queueName := config.GetString(ctx, "sqs.name")
	//
	//queue, err := sqs.NewStandardQueue(ctx, queueName, &sqs.Config{
	//	Account:  config.GetString(ctx, "sqs.account"),
	//	Endpoint: config.GetString(ctx, "sqs.endpoint"),
	//	Region:   config.GetString(ctx, "sqs.region"),
	//})
	//

	queue := psqs.QueueGlobal

	if err != nil || queue == nil {
		cusErr = error2.NewCustomError(error3.SqsInitializeErr, err.Error())
		return
	}

	publisher := sqs.NewPublisher(queue)
	m := &sqs.Message{
		Value: msg,
	}

	err = publisher.Publish(ctx, m)
	if err != nil {
		cusErr = error2.NewCustomError(error3.SqsPublishErr, err.Error())
		return
	}

	fmt.Println(101)

	log.Infof("Email send message is successfully pushed to SQS | %v", req)
	return
}
