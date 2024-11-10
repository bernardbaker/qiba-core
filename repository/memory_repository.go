package repository

import (
	"github.com/bernardbaker/qiba.core/domain"
)

type MemoryMessageRepository struct {
	messages []*domain.Message
}

func NewMemoryMessageRepository() *MemoryMessageRepository {
	return &MemoryMessageRepository{
		messages: make([]*domain.Message, 0),
	}
}

func (r *MemoryMessageRepository) Save(msg *domain.Message) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *MemoryMessageRepository) SaveMessage(msg *domain.Message) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *MemoryMessageRepository) GetMessages(viewerID string) ([]*domain.Message, error) {
	var result []*domain.Message
	for _, msg := range r.messages {
		if msg.Viewer == viewerID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (r *MemoryMessageRepository) GetAll() ([]*domain.Message, error) {
	var result []*domain.Message
	result = append(result, r.messages...)
	return result, nil
}
