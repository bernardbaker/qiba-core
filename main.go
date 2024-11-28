package main

import (
	"log"
	"net"
	"os"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
	"github.com/bernardbaker/qiba.core/ports"
	"github.com/bernardbaker/qiba.core/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Repository interfaces
type RepositoryType string

const (
	InMemory RepositoryType = "inmemory"
	MongoDB  RepositoryType = "mongodb"
)

func getRepositories(repoType RepositoryType) (
	gameRepo ports.GameRepository,
	userRepo ports.UserRepository,
	leaderboardRepo ports.LeaderboardRepository,
	referralRepo ports.ReferralRepository,
) {
	switch repoType {
	case InMemory:
		return infrastructure.NewInMemoryGameRepository(),
			infrastructure.NewInMemoryUserRepository(),
			infrastructure.NewInMemoryLeaderboardRepository(),
			infrastructure.NewInMemoryReferralRepository()
	// Add cases for other repository types
	case MongoDB:
		return infrastructure.NewMongoDbGameRepository(),
			infrastructure.NewMongoDbUserRepository(),
			infrastructure.NewMongoDbLeaderboardRepository(),
			infrastructure.NewMongoDbReferralRepository()
	default:
		log.Printf("Unknown repository type %s, falling back to in-memory", repoType)
		return infrastructure.NewInMemoryGameRepository(),
			infrastructure.NewInMemoryUserRepository(),
			infrastructure.NewInMemoryLeaderboardRepository(),
			infrastructure.NewInMemoryReferralRepository()
	}
}

func main() {
	// Get repository type from environment variable
	repoType := RepositoryType(os.Getenv("REPOSITORY_TYPE"))
	if repoType == "" {
		repoType = InMemory
		log.Printf("No repository type specified, defaulting to %s", repoType)
	}
	// Get port number from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}
	// Initialize repositories based on type
	gameRepo, userRepo, leaderboardRepo, referralRepo := getRepositories(repoType)

	// Initialize encrypter
	encrypter := infrastructure.NewEncrypter([]byte("mysecretencryptionkey1234567890a"))
	// Initialize game service
	service := app.NewGameService(gameRepo, userRepo, leaderboardRepo, encrypter)
	// Initialize referral service
	referralService := app.NewReferralService(referralRepo)

	// Initialize the leader board
	if os.Getenv("ENV") == "development" {
		service.CreateLeaderboard("qiba", true)
	} else {
		service.CreateLeaderboard("qiba", true)
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
