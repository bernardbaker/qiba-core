package domain

// Message represents a chat message
type Message struct {
	ID      string
	Content string
	Viewer  string
}

type MessageRepository interface {
	Save(message *Message) error
	GetAll() ([]*Message, error)
}

// Publisher defines an interface for publishing messages
type Publisher interface {
	Publish(string, []string) error
}

// Receiver defines an interface for receiving messages
type Receiver interface {
	Receive() ([]string, error)
}
