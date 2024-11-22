package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/infrastructure"
	"github.com/bernardbaker/qiba.core/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	_ "google.golang.org/grpc/encoding/gzip"
)

func main() {
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
	service.CreateLeaderboard("qiba")

	var serverOpts []grpc.ServerOption

	// Check if we're in development mode
	if os.Getenv("ENV") == "development" {
		// Add your production TLS configuration here
		// Example: Load certificates from files or secrets management
		// creds := credentials.NewTLS(&tls.Config{...})
		// serverOpts = append(serverOpts, grpc.Creds(creds))
	} else {
		// TLS with InsecureSkipVerify for local development only
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false,
		})
		serverOpts = append(
			serverOpts,
			grpc.Creds(creds),
			grpc.MaxRecvMsgSize(2000*1024*1024),
			grpc.MaxSendMsgSize(2000*1024*1024),
		)
	}

	// Setup and start gRPC server
	server := grpc.NewServer(serverOpts...)
	// Register game service
	proto.RegisterGameServiceServer(server, infrastructure.NewGameServer(service))
	// Register referral service
	proto.RegisterReferralServiceServer(server, infrastructure.NewReferralServer(referralService, service))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Setting new Logger
	grpcLog := grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr)
	grpclog.SetLoggerV2(grpcLog)

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server listening at %v", listener.Addr())

	log.Printf("Starting gRPC server on port %s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
