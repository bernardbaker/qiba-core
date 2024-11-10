package repository

import (
	"github.com/bernardbaker/qiba.core/ports"
)

type MemoryMessageRepository struct {
	messages []*ports.Message
}

func NewMemoryMessageRepository() *MemoryMessageRepository {
	return &MemoryMessageRepository{
		messages: make([]*ports.Message, 0),
	}
}

func (r *MemoryMessageRepository) Save(msg *ports.Message) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *MemoryMessageRepository) SaveMessage(msg *ports.Message) error {
	r.messages = append(r.messages, msg)
	return nil
}

func (r *MemoryMessageRepository) GetMessages(viewerID string) ([]*ports.Message, error) {
	var result []*ports.Message
	for _, msg := range r.messages {
		if msg.Viewer == viewerID {
			result = append(result, msg)
		}
	}
	return result, nil
}

func (r *MemoryMessageRepository) GetAll() ([]*ports.Message, error) {
	var result []*ports.Message
	result = append(result, r.messages...)
	return result, nil
}
