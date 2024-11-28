package infrastructure

import (
	"errors"
	"sync"

	"github.com/bernardbaker/qiba.core/domain"
)

type MongoDbGameRepository struct {
	games map[string]*domain.Game
	mutex sync.RWMutex
}

func NewMongoDbGameRepository() *MongoDbGameRepository {
	return &MongoDbGameRepository{
		games: make(map[string]*domain.Game),
	}
}

// SaveGame stores a new game in the in-memory map
func (repo *MongoDbGameRepository) SaveGame(game *domain.Game) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.games[game.ID] = game
	return nil
}

// GetGame retrieves a game by its ID
func (repo *MongoDbGameRepository) GetGame(gameID string) (*domain.Game, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	game, exists := repo.games[gameID]
	if !exists {
		return nil, errors.New("game not found")
	}
	return game, nil
}

// UpdateGame updates an existing game in the in-memory map
func (repo *MongoDbGameRepository) UpdateGame(game *domain.Game) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.games[game.ID]
	if !exists {
		return errors.New("game not found")
	}
	repo.games[game.ID] = game
	return nil
}

// GetGamesByUser
func (repo *MongoDbGameRepository) GetGamesByUser(userID string) ([]*domain.Game, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	var games []*domain.Game
	for _, game := range repo.games {
		if game.UserID == userID {
			games = append(games, game)
		}
	}
	return games, nil
}
