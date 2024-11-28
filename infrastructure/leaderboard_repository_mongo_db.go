package infrastructure

import (
	"errors"
	"sync"

	"github.com/bernardbaker/qiba.core/domain"
)

type MongoDbLeaderboardRepository struct {
	store map[string]*domain.Table
	mutex sync.RWMutex
}

func NewMongoDbLeaderboardRepository() *MongoDbLeaderboardRepository {
	return &MongoDbLeaderboardRepository{
		store: make(map[string]*domain.Table),
	}
}

// SaveGame stores a new leaderboard in the in-memory map
func (repo *MongoDbLeaderboardRepository) SaveLeaderboard(table *domain.Table) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.store[table.ID] = table
	return nil
}

// GetGame retrieves a table by its ID
func (repo *MongoDbLeaderboardRepository) GetLeaderboard(tableID string) (*domain.Table, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	table, exists := repo.store[tableID]
	if !exists {
		return nil, errors.New("table not found")
	}
	return table, nil
}

// UpdateGame updates an existing table in the in-memory map
func (repo *MongoDbLeaderboardRepository) UpdateLeaderboard(table *domain.Table) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.store[table.ID]
	if !exists {
		return errors.New("table not found")
	}
	repo.store[table.ID] = table
	return nil
}

// UpdateGame updates an existing table in the in-memory map
func (repo *MongoDbLeaderboardRepository) AddEntryToLeaderboard(table *domain.Table, entry *domain.GameEntry) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.store[table.ID]
	if !exists {
		return errors.New("table not found")
	}
	table.Entries = append(table.Entries, *entry)
	repo.store[table.ID] = table
	return nil
}
