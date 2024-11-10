package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChatService defines the chat service operations
type ChatService struct {
	repo                                 domain.MessageRepository
	publisher                            domain.Publisher
	receiver                             domain.Receiver
	proto.UnimplementedChatServiceServer // Embed this to get default implementations
}

// NewChatService creates a new ChatService
func NewChatService(repo domain.MessageRepository, publisher domain.Publisher, receiver domain.Receiver) *ChatService {
	return &ChatService{
		repo:      repo,
		publisher: publisher,
		receiver:  receiver,
	}
}

// SendMessage saves and publishes a message from a viewer
func (s *ChatService) SendMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	message := domain.Message{
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
	// err = s.publisher.Publish(message.Content, []string{message.Viewer})
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to publish message: %v", err)
	// }

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
		message := domain.Message{
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
	// err := s.publisher.Publish(content, viewers)
	// if err != nil {
	// 	return nil, status.Errorf(codes.Internal, "failed to publish broadcast message: %v", err)
	// }

	return &proto.BroadcastResponse{
		Status: "broadcast sent",
	}, nil
}

func (s *ChatService) ReceiveMessages() error {
	// Receive messages from memory
	messages, err := s.repo.GetAll()
	if err != nil {
		return fmt.Errorf("failed to get all messages: %w", err)
	}

	for _, message := range messages {
		// Print the message content
		fmt.Printf("Received message: %s\n", message.Content)
	}

	return nil
}

// generateMessageID is a helper function to generate unique message IDs
func generateMessageID() string {
	// Simple unique ID generation for demonstration purposes
	// You might want to use a more robust mechanism (like UUIDs) in production
	return fmt.Sprintf("msg-%d", time.Now().UnixNano())
}
