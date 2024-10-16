package grpc

import (
	"context"

	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/domain"
	// Generated from proto// Proto generated code
)

type GRPCChatService struct {
	chatService domain.ChatService
}

func NewGRPCChatService(svc domain.ChatService) *GRPCChatService {
	return &GRPCChatService{svc}
}

func (s *GRPCChatService) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	message, err := s.chatService.SendMessage(req.ViewerId, req.Content)
	if err != nil {
		return nil, err
	}
	return &pb.MessageResponse{
		MessageId: message.ID,
		Status:    "sent",
	}, nil
}

func (s *GRPCChatService) SendBroadcast(ctx context.Context, req *pb.BroadcastRequest) (*pb.BroadcastResponse, error) {
	err := s.chatService.BroadcastMessage(req.Message, req.Viewers)
	if err != nil {
		return nil, err
	}
	return &pb.BroadcastResponse{Status: "broadcasted"}, nil
}
