package domain

import (
	"time"

	"github.com/google/uuid"

	"math/rand"
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
	// game.generateObjectSequence()
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
