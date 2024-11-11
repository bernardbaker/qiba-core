package main

import (
	"log"
	"net"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
	"github.com/bernardbaker/qiba.core/proto"

	"google.golang.org/grpc"
)

func main() {
	// Initialize repository, encrypter, and game service
	repo := infrastructure.NewInMemoryGameRepository() // Our in-memory game repository
	encrypter := infrastructure.NewEncrypter([]byte("mysecretencryptionkey1234567890!"))
	service := app.NewGameService(repo, encrypter)

	// Setup and start gRPC server
	server := grpc.NewServer()
	proto.RegisterGameServiceServer(server, infrastructure.NewGameServer(service))

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting gRPC server on port 50051...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
