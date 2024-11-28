package mocks

import (
	"math/rand"
	"time"

	"github.com/bernardbaker/qiba.core/domain"
)

func GenerateMockData(entries int) []*domain.GameEntry {
	usernames := []string{"rogue", "mage", "warrior", "healer", "hunter", "rogue", "mage", "paladin", "extralongusernamefromsomewhere"}
	mockData := make([]*domain.GameEntry, 0, entries)

	// Seed random number generator
	rand.NewSource(time.Now().UnixNano())

	for i := 0; i < entries; i++ {
		// Random username from the list
		username := usernames[rand.Intn(len(usernames))]

		// Random score between 0 and 10
		score := rand.Int31n(10)

		// Random timestamp (allow some duplicates)
		baseTime := time.Now().Add(-time.Duration(rand.Intn(100)) * time.Minute)
		timestamp := baseTime.Add(time.Duration(rand.Intn(60)) * time.Second)

		entry := &domain.GameEntry{
			User:      domain.User{Username: username},
			Score:     score,
			Timestamp: timestamp,
		}
		mockData = append(mockData, entry)
	}

	return mockData
}
