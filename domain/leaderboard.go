package domain

import (
	"time"
)

type GameEntry struct {
	User      User
	Score     int32
	Timestamp time.Time
}

type Table struct {
	ID      string
	Entries []GameEntry
}

type LeaderboardEntry struct {
	Username  string
	Score     int32
	Timestamp time.Time
}

func NewLeaderboard(name string) *Table {
	board := &Table{
		ID:      name,
		Entries: []GameEntry{},
	}
	return board
}

func NewLeaderboardObject(user User, score int32) *GameEntry {
	entry := &GameEntry{
		User:      user,
		Score:     score,
		Timestamp: time.Now(),
	}
	return entry
}
