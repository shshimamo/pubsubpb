package main

import (
	"encoding/base64"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/protobuf/proto"
	"github.com/shshimamo/pubsubpb/pb"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("ap-northeast-1"),
		Endpoint: aws.String("http://localhost:4566"),
	}))

	svc := sqs.New(sess)

	//queueUrl := "http://localhost:4566/000000000000/message-queue"
	queueUrl := "http://sqs.ap-northeast-1.localhost.localstack.cloud:4566/000000000000/message-queue"

	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: aws.Int64(1),
		WaitTimeSeconds:     aws.Int64(10),
	})
	if err != nil {
		log.Fatalf("Failed to receive message: %v", err)
	}

	if len(result.Messages) == 0 {
		log.Println("No messages received")
		return
	}

	// Decode the Base64 encoded message
	decodedPayload, err := base64.StdEncoding.DecodeString(*result.Messages[0].Body)
	if err != nil {
		log.Fatalf("Failed to decode message: %v", err)
	}

	wrappedMsg := &pb.WrappedMessage{}
	err = proto.Unmarshal(decodedPayload, wrappedMsg)
	if err != nil {
		log.Fatalf("Failed to unmarshal wrapped message: %v", err)
	}

	switch wrappedMsg.Type {
	case pb.WrappedMessage_TYPE_A:
		msgA := &pb.MessageTypeA{}
		err := proto.Unmarshal(wrappedMsg.Payload, msgA)
		if err != nil {
			log.Fatalf("Failed to unmarshal message type A: %v", err)
		}
		log.Printf("Received MessageTypeA: %v", msgA)
	case pb.WrappedMessage_TYPE_B:
		msgB := &pb.MessageTypeB{}
		err := proto.Unmarshal(wrappedMsg.Payload, msgB)
		if err != nil {
			log.Fatalf("Failed to unmarshal message type B: %v", err)
		}
		log.Printf("Received MessageTypeB: %v", msgB)
	default:
		log.Fatalf("Unknown message type: %v", wrappedMsg.Type)
	}
}
