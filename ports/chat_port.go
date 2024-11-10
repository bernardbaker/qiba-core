package ports

import (
	"context"

	"github.com/bernardbaker/qiba.core/proto"
)

type ChatPort interface {
	SendMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error)
	SendBroadcast(ctx context.Context, req *proto.BroadcastRequest) (*proto.BroadcastResponse, error)
	ReceiveMessages() error
}
