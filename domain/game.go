package domain

import (
	"time"

	"github.com/google/uuid"

	"math/rand"
)

type Game struct {
	ID        string       `bson:"ID"`
	Score     int32        `bson:"Score"`
	ObjectSeq []GameObject `bson:"ObjectSeq"`
	StartTime time.Time    `bson:"StartTime"`
	EndTime   time.Time    `bson:"EndTime"`
	UserID    string       `bson:"UserID"`
}

type GameObject struct {
	ID        string    `bson:"ID"`
	Type      string    `bson:"Type"`
	Timestamp time.Time `bson:"Timestamp"`
}

// Generate a new game with random object sequence
func NewGame(userId string) *Game {
	game := &Game{
		ID:        uuid.New().String(),
		StartTime: time.Now(),
		UserID:    userId,
		EndTime:   time.Now(),
	}
	return game
}

// Generates a random sequence of objects for the game
func (g *Game) GenerateObjectSequence() {
	// Populate with 10 sample objects

	isTypeA := rand.Intn(2) == 0

	g.ObjectSeq = append(g.ObjectSeq, GameObject{
		ID:        uuid.New().String(),
		Type:      map[bool]string{true: "a", false: "b"}[isTypeA],
		Timestamp: time.Now(),
	})

}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, t.Location())
}
