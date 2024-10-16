package repository

import (
	"github.com/bernardbaker/streamlit.chat.using.hexagonal.pattern/domain"
)

type MemoryMessageRepository struct {
	messages []*domain.Message
}

func NewMemoryMessageRepository() *MemoryMessageRepository {
	return &MemoryMessageRepository{}
}

func (r *MemoryMessageRepository) SaveMessage(msg *domain.Message) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *MemoryMessageRepository) GetMessages(viewerID string) ([]*domain.Message, error) {
	var result []*domain.Message
	for _, msg := range r.messages {
		if msg.ViewerID == viewerID {
			result = append(result, msg)
		}
	}
	return result, nil
}
