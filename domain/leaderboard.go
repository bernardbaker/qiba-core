package domain

import (
	"cmp"
	"slices"
	"time"
)

type GameEntry struct {
	User      User      `bson:"User"`
	Score     int32     `bson:"Score"`
	Timestamp time.Time `bson:"Timestamp"`
}

type Table struct {
	ID      string      `bson:"id"`
	Entries []GameEntry `bson:"entries"`
}

type LeaderboardEntry struct {
	Username string
	Score    int32
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

type UserScore struct {
	DisplayName string
	TotalScore  int32
	User        User
	// LastPlayed  time.Time
}

func GroupAndTotalScoresByUser(table *Table) map[string]int32 {
	// Create a map to store total scores for each unique user
	userTotals := make(map[string]int32)

	// Iterate through entries and sum scores by user
	for _, entry := range table.Entries {
		displayName := getUserDisplayName(entry.User)
		userTotals[displayName] += entry.Score
	}

	return userTotals
}

// Helper function to get the display name based on priority
func getUserDisplayName(user User) string {
	if user.Username != "" {
		return user.Username
	} else if user.FirstName != "" {
		return user.FirstName
	}
	return user.LastName
}

func GroupAndTotalScoresByUserSorted(table *Table) []UserScore {
	// Create maps to store totals and track last played time
	userTotals := make(map[string]int32)
	userDetails := make(map[string]User)

	// Calculate totals and track most recent timestamp for each user
	for _, entry := range table.Entries {
		displayName := getUserDisplayName(entry.User)
		userTotals[displayName] += entry.Score
		userDetails[displayName] = entry.User
	}

	// Convert map to slice for sorting
	var results []UserScore
	for displayName, totalScore := range userTotals {
		results = append(results, UserScore{
			DisplayName: displayName,
			TotalScore:  totalScore,
			User:        userDetails[displayName],
		})
	}

	// Sort by total score (highest first) and then by last played time
	slices.SortFunc(results, func(a, b UserScore) int {
		return cmp.Compare(b.TotalScore, a.TotalScore)
	})

	return results
}
