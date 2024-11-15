package main

import (
	"log"
	"net"
	"os"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
	"github.com/bernardbaker/qiba.core/proto"

	"google.golang.org/grpc"
)

func main() {
	// Initialize repository, encrypter, and game service
	repo := infrastructure.NewInMemoryGameRepository()     // Our in-memory game repository
	userRepo := infrastructure.NewInMemoryUserRepository() // Our in-memory game repository
	encrypter := infrastructure.NewEncrypter([]byte("mysecretencryptionkey1234567890a"))
	service := app.NewGameService(repo, userRepo, encrypter)

	// Initialize the referral repository and referral service
	referralRepo := infrastructure.NewInMemoryReferralRepository()
	referralService := app.NewReferralService(referralRepo)

	// Setup and start gRPC server
	server := grpc.NewServer()
	// Register game service
	proto.RegisterGameServiceServer(server, infrastructure.NewGameServer(service))
	// Register referral service
	proto.RegisterReferralServiceServer(server, infrastructure.NewReferralServer(referralService, service))

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
		log.Printf("defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Starting gRPC server on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
