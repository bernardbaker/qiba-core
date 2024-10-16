package main

import (
	"log"
	"net"

	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/aws"
	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/domain"
	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/grpc"
	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/repository"

	pb "github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/proto"
	"google.golang.org/grpc"
)

func main() {
	// Create the repository
	repo := repository.NewMemoryMessageRepository()

	// Create SNS Publisher and SQS Receiver
	snsPublisher := aws.NewSNSPublisher("sns-topic-arn")
	sqsReceiver := aws.NewSQSReceiver("sqs-queue-url")

	// Create the domain service
	chatService := domain.NewChatService(repo, snsPublisher, sqsReceiver)

	// Set up the gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	chatGrpcService := grpc.NewGRPCChatService(chatService)
	pb.RegisterChatServiceServer(grpcServer, chatGrpcService)

	log.Printf("Server listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
