package domain

import (
	"cmp"
	"slices"
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

func OrderLeaderboard(table *Table) *Table {
	// Sort table by Score and Timestamp using slices.SortFunc.
	// The newest score is further up the leaderboard.
	compareFunc := func(a, b GameEntry) int {
		if a.Score != b.Score {
			return cmp.Compare(b.Score, a.Score)
		}
		// If scores are equal, compare timestamps
		return cmp.Compare(b.Timestamp.UnixMilli(), a.Timestamp.UnixMilli())
	}

	slices.SortFunc(table.Entries, compareFunc)
	return table
}
