package domain

import (
	"fmt"
	"time"

	"context"

	proto "github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Message represents a chat message
type Message struct {
	ID      string
	Content string
	Viewer  string
}

type MessageRepository interface {
	Save(message *Message) error
}

// Publisher defines an interface for publishing messages
type Publisher interface {
	Publish(string, []string) error
}

// Receiver defines an interface for receiving messages
type Receiver interface {
	Receive() ([]string, error)
}

// ChatService defines the chat service operations
type ChatService struct {
	repo                                 MessageRepository
	publisher                            Publisher
	receiver                             Receiver
	proto.UnimplementedChatServiceServer // Embed this to get default implementations
}

// NewChatService creates a new ChatService
func NewChatService(repo MessageRepository, publisher Publisher, receiver Receiver) *ChatService {
	return &ChatService{
		repo:      repo,
		publisher: publisher,
		receiver:  receiver,
	}
}

// SendMessage saves and publishes a message from a viewer
func (s *ChatService) SendMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	message := Message{
		ID:      generateMessageID(),
		Content: req.Content,
		Viewer:  req.ViewerId,
	}

	// Save the message to the repository
	err := s.repo.Save(&message)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save message: %v", err)
	}

	// Publish the message
	err = s.publisher.Publish(message.Content, []string{message.Viewer})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	}

	return &proto.MessageResponse{
		MessageId: message.ID,
		Status:    "sent",
	}, nil
}

// SendBroadcast publishes a message to multiple viewers
func (s *ChatService) SendBroadcast(ctx context.Context, req *proto.BroadcastRequest) (*proto.BroadcastResponse, error) {
	// Extract content and viewers from the request
	content := req.Message
	viewers := req.Viewers

	for _, viewerID := range viewers {
		message := Message{
			ID:      generateMessageID(),
			Content: content,
			Viewer:  viewerID,
		}

		// Save the message to the repository
		err := s.repo.Save(&message)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to save broadcast message: %v", err)
		}
	}

	// Publish the message to all viewers
	err := s.publisher.Publish(content, viewers)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to publish broadcast message: %v", err)
	}

	return &proto.BroadcastResponse{
		Status: "broadcast sent",
	}, nil
}

// ReceiveMessages retrieves and processes messages from the SQS queue
func (s *ChatService) ReceiveMessages() error {
	// Receive messages from SQS
	messages, err := s.receiver.Receive()
	if err != nil {
		return fmt.Errorf("failed to receive messages: %w", err)
	}

	// Process each received message
	for _, content := range messages {
		message := Message{
			ID:      generateMessageID(),
			Content: content,
			Viewer:  "system", // Set the viewer ID for system messages or adjust based on logic
		}

		// Save the message to the repository
		err := s.repo.Save(&message)
		if err != nil {
			return fmt.Errorf("failed to save received message: %w", err)
		}

		// Optionally, publish the received message if required
		err = s.publisher.Publish(message.Content, []string{message.Viewer})
		if err != nil {
			return fmt.Errorf("failed to publish received message: %w", err)
		}
	}

	return nil
}

// generateMessageID is a helper function to generate unique message IDs
func generateMessageID() string {
	// Simple unique ID generation for demonstration purposes
	// You might want to use a more robust mechanism (like UUIDs) in production
	return fmt.Sprintf("msg-%d", time.Now().UnixNano())
}
