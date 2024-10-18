package main

import (
	"context"
	"testing"
	"time"

	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestSendMessage(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new ChatService client
	client := proto.NewChatServiceClient(conn)

	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Prepare the message request
	req := &proto.MessageRequest{
		Content:  "Hello, gRPC!",
		ViewerId: "testUser123",
	}

	// Send the message
	resp, err := client.SendMessage(ctx, req)
	if err != nil {
		t.Fatalf("could not send message: %v", err)
	}

	// Check the response
	if resp.MessageId == "" {
		t.Errorf("expected non-empty message ID, got empty string")
	}
	if resp.Status != "sent" {
		t.Errorf("expected status 'sent', got %s", resp.Status)
	}

	t.Logf("Message sent successfully. Message ID: %s, Status: %s", resp.MessageId, resp.Status)
}

func TestBroadcastMessage(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Create a new ChatService client
	client := proto.NewChatServiceClient(conn)

	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Prepare the message request
	req := &proto.BroadcastRequest{
		Message: "Hello, gRPC!",
		Viewers: []string{"testUser1", "testUser2"},
	}

	// Send the message
	resp, err := client.SendBroadcast(ctx, req)
	if err != nil {
		t.Fatalf("could not broadcast message: %v", err)
	}

	// Check the response
	if resp.Status != "broadcast sent" {
		t.Errorf("expected status 'sent', got %s", resp.Status)
	}

	t.Logf("Message sent successfully. Status: %s", resp.Status)
}
