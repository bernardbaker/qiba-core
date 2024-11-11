package domain

import (
	"time"

	"github.com/google/uuid"
)

type Game struct {
	ID        string
	Score     int32
	ObjectSeq []GameObject
	StartTime time.Time
	EndTime   time.Time
}

type GameObject struct {
	ID        string
	Type      string // "a" for add point, "b" for subtract point
	Timestamp time.Time
}

// Generate a new game with random object sequence
func NewGame() *Game {
	game := &Game{
		ID:        uuid.New().String(),
		StartTime: time.Now(),
	}
	game.generateObjectSequence()
	return game
}

// Generates a random sequence of objects for the game
func (g *Game) generateObjectSequence() {
	// Populate with 10 sample objects
	for i := 0; i < 10; i++ {
		g.ObjectSeq = append(g.ObjectSeq, GameObject{
			ID:        uuid.New().String(),
			Type:      map[bool]string{true: "a", false: "b"}[i%2 == 0],
			Timestamp: time.Now().Add(time.Duration(i) * time.Second),
		})
	}
}
