package domain

import "time"

type Message struct {
    ID        string
    ViewerID  string
    Content   string
    Timestamp int64
}

type ChatService interface {
    SendMessage(viewerID string, content string) (*Message, error)
    BroadcastMessage(content string, viewerIDs []string) error
}

type MessageRepository interface {
    SaveMessage(msg *Message) error
    GetMessages(viewerID string) ([]*Message, error)
}

func NewMessage(id, viewerID, content string) *Message {
    return &Message{
        ID:        id,
        ViewerID:  viewerID,
        Content:   content,
        Timestamp: time.Now().Unix(),
    }
}
