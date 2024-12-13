package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
	"github.com/bernardbaker/qiba.core/ports"
	"github.com/bernardbaker/qiba.core/proto"

	"google.golang.org/grpc"
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
	// case MongoDB:
	// 	return infrastructure.NewInMemoryGameRepository(),
	// 		infrastructure.NewInMemoryUserRepository(),
	// 		infrastructure.NewInMemoryLeaderboardRepository(),
	// 		infrastructure.NewInMemoryReferralRepository()
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

func isPortInUse(host string, port string) bool {
	address := net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return false
	}
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}

func disableAllLogs() {
	// Redirect stderr to null device
	null, _ := os.Open(os.DevNull)
	os.Stderr = null

	// Redirect stdout to null device if needed
	os.Stdout = null
}

func main() {
	if os.Getenv("ENV") != "development" {
		disableAllLogs()
	}

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

	// Check if port is already in use
	if isPortInUse("0.0.0.0", port) {
		log.Printf("Port %s is already in use", port)
		return
	}

	// Initialize repositories based on type
	gameRepo, userRepo, leaderboardRepo, referralRepo := getRepositories(repoType)

	// Initialize encrypter
	encrypter := infrastructure.NewEncrypter([]byte("mysecretencryptionkey1234567890a"))
	// Initialize game service
	service := app.NewGameService(gameRepo, userRepo, leaderboardRepo, encrypter)
	// Initialize referral service
	referralService := app.NewReferralService(referralRepo)

	// Prepopulate the leaderboard
	// TODO: if the users score is not in the top 100 find it and display it.
	prepopulate := false
	if prepopulate {
		// Initialize the leader board
		service.CreateLeaderboard("dev", prepopulate)
	} else {
		service.CreateLeaderboard("qiba", false)
	}

	// Setting new Logger
	// grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	// grpclog.SetLoggerV2(grpcLog)

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Server listening at %v", listener.Addr().String())

	server := grpc.NewServer()

	// Register gRPC services
	proto.RegisterGameServiceServer(server, infrastructure.NewGameServer(service))
	proto.RegisterReferralServiceServer(server, infrastructure.NewReferralServer(referralService, service))

	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	log.Printf("Starting gRPC server on port %s", port)
}
