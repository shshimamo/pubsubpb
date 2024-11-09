package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/protobuf/proto"
	"github.com/shshimamo/pubsubpb/pb"
)

func main() {
	createSQS()

	msgA := &pb.MessageTypeA{
		Field1: "example",
		Field2: 123,
	}

	payload, err := proto.Marshal(msgA)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	wrappedMsg := &pb.WrappedMessage{
		Type:    pb.WrappedMessage_TYPE_A,
		Payload: payload,
	}

	wrappedPayload, err := proto.Marshal(wrappedMsg)
	if err != nil {
		log.Fatalf("Failed to marshal wrapped message: %v", err)
	}

	// Encode wrappedPayload to Base64
	encodedPayload := base64.StdEncoding.EncodeToString(wrappedPayload)
	fmt.Printf("encodedPayload: %v\n", encodedPayload)

	// Send wrappedPayload to SNS or SQS
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}))

	svc := sqs.New(sess)

	queueUrl := "http://sqs.ap-northeast-1.localhost.localstack.cloud:4566/000000000000/message-queue"

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: aws.String(encodedPayload),
	})
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	log.Println("Successfully sent message")
}

func createSQS() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}))

	svc := sqs.New(sess)

	result, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String("message-queue"),
	})
	if err != nil {
		log.Fatalf("Unable to create queue: %v", err)
	}

	log.Printf("Successfully created queue %s", *result.QueueUrl)
}
