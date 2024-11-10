package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/bernardbaker/qiba.core/aws"
	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/proto"
	"github.com/bernardbaker/qiba.core/repository"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ChatServer implements the gRPC server for chat service
type ChatServer struct {
	domain.ChatService
}

// NewChatServer creates a new instance of ChatServer
func NewChatServer(chatService *domain.ChatService) *ChatServer {
	return &ChatServer{
		ChatService: *chatService,
	}
}

// Implement the gRPC service methods (SendMessage, BroadcastMessage, etc.)

func main() {
	// Initialize repository (e.g., in-memory or database-backed)
	repo := repository.NewMemoryMessageRepository() // Assuming in-memory for now

	// Initialize AWS SNS Publisher with an example SNS Topic ARN
	snsPublisher := aws.NewSNSPublisher("sns-topic-arn")

	// Initialize AWS SQS Receiver with an example SQS Queue URL
	sqsReceiver := aws.NewSQSReceiver("sqs-queue-url")

	// Create the ChatService domain layer
	chatService := domain.NewChatService(repo, snsPublisher, sqsReceiver)

	// Initialize the gRPC server
	grpcServer := grpc.NewServer()

	// Register your gRPC service implementation
	chatServer := NewChatServer(chatService)
	// Assuming RegisterChatServiceServer is the function to register the gRPC service, replace with the actual generated code --replaces
	// Receive() ([]string, error)
	proto.RegisterChatServiceServer(grpcServer, chatServer)

	// Enable gRPC reflection (useful for debugging)
	reflection.Register(grpcServer)

	// Set up the gRPC server listener on a specific port
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	// Start a background goroutine to periodically poll and process messages from SQS
	go func() {
		for {
			// Receive messages from the SQS queue
			err := chatService.ReceiveMessages()
			if err != nil {
				log.Printf("Error receiving messages from Memory/SQS: %v", err)
			}

			// Sleep for 10 seconds before the next polling attempt (adjust the interval as needed)
			time.Sleep(10 * time.Second)
		}
	}()

	// Start the gRPC server
	fmt.Println("Starting gRPC server on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve gRPC server: %v", err)
	}
}
