package main

import (
	"log"
	"net"
	"os"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
	"github.com/bernardbaker/qiba.core/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	// Initialize repository, encrypter, and game service
	repo := infrastructure.NewInMemoryGameRepository()                   // Our in-memory game repository
	userRepo := infrastructure.NewInMemoryUserRepository()               // Our in-memory game repository
	leaderboardRepo := infrastructure.NewInMemoryLeaderboardRepository() // Our in-memory game repository
	encrypter := infrastructure.NewEncrypter([]byte("mysecretencryptionkey1234567890a"))
	service := app.NewGameService(repo, userRepo, leaderboardRepo, encrypter)
	// Initialize the referral repository and referral service
	referralRepo := infrastructure.NewInMemoryReferralRepository()
	referralService := app.NewReferralService(referralRepo)
	// Initialize the leader board

	if os.Getenv("ENV") == "development" {
		service.CreateLeaderboard("qiba", true)
	} else {
		service.CreateLeaderboard("qiba", false)
	}

	// Setting new Logger
	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", listener.Addr().String())

	var server *grpc.Server
	if os.Getenv("ENV") == "development" {
		server = grpc.NewServer()
	} else {
		// fallbackCreds := credentials.NewTLS(&tls.Config{
		// 	InsecureSkipVerify: true, // Skip certificate verification for testing purposes
		// })
		// Use TLS credentials or custom options
		// creds, err := xds.NewServerCredentials(xds.ServerOptions{
		// 	FallbackCreds: fallbackCreds,
		// })
		// if err != nil {
		// 	log.Fatalf("failed to create server credentials: %v", err)
		// }
		// server = grpc.NewServer(grpc.Creds(fallbackCreds))
		server = grpc.NewServer()
	}

	// Register gRPC services
	proto.RegisterGameServiceServer(server, infrastructure.NewGameServer(service))
	proto.RegisterReferralServiceServer(server, infrastructure.NewReferralServer(referralService, service))

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("Starting gRPC server on port %s", port)
}
