package app

import (
	"errors"
	"testing"
	"time"

	"github.com/bernardbaker/qiba.core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repository
type MockGameRepository struct {
	mock.Mock
}

func (m *MockGameRepository) SaveGame(game *domain.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

func (m *MockGameRepository) GetGame(id string) (*domain.Game, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Game), args.Error(1)
}

func (m *MockGameRepository) UpdateGame(game *domain.Game) error {
	args := m.Called(game)
	return args.Error(0)
}

// Mock Encrypter
type MockEncrypter struct {
	mock.Mock
}

func (m *MockEncrypter) EncryptGameData(objects interface{}) (string, string, error) {
	args := m.Called(objects)
	return args.String(0), args.String(1), args.Error(2)
}

func TestNewGameService(t *testing.T) {
	repo := new(MockGameRepository)
	encrypter := new(MockEncrypter)

	service := NewGameService(repo, encrypter)

	assert.NotNil(t, service)
	assert.Equal(t, repo, service.repo)
	assert.Equal(t, encrypter, service.encrypter)
}

func TestStartGame(t *testing.T) {
	t.Run("successful game start", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		repo.On("SaveGame", mock.AnythingOfType("*domain.Game")).Return(nil)
		encrypter.On("EncryptGameData", mock.AnythingOfType("[]domain.GameObject")).
			Return("encrypted_data", "hmac_value", nil)

		encryptedData, hmac, _, err := service.StartGame()

		assert.NoError(t, err)
		assert.Equal(t, "encrypted_data", encryptedData)
		assert.Equal(t, "hmac_value", hmac)
		repo.AssertExpectations(t)
		encrypter.AssertExpectations(t)
	})

	t.Run("save game fails", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		repo.On("SaveGame", mock.AnythingOfType("*domain.Game")).
			Return(assert.AnError)

		encryptedData, hmac, gameId, err := service.StartGame()

		assert.Error(t, err)
		assert.Empty(t, encryptedData)
		assert.Empty(t, hmac)
		assert.Empty(t, gameId)
		repo.AssertExpectations(t)
	})

	t.Run("encryption fails", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		repo.On("SaveGame", mock.AnythingOfType("*domain.Game")).Return(nil)
		encrypter.On("EncryptGameData", mock.AnythingOfType("[]domain.GameObject")).
			Return("", "", assert.AnError)

		encryptedData, hmac, _, err := service.StartGame()

		assert.Error(t, err)
		assert.Empty(t, encryptedData)
		assert.Empty(t, hmac)
		repo.AssertExpectations(t)
		encrypter.AssertExpectations(t)
	})
}

func TestTap(t *testing.T) {
	t.Run("successful tap on type 'a'", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		game := &domain.Game{
			ObjectSeq: []domain.GameObject{
				{
					ID:        "obj1",
					Type:      "a",
					Timestamp: time.Now().Add(-time.Minute),
				},
			},
			Score: 0,
		}

		repo.On("GetGame", "game1").Return(game, nil)
		repo.On("UpdateGame", mock.AnythingOfType("*domain.Game")).Return(nil)

		success, err := service.Tap("game1", "obj1", time.Now())

		assert.NoError(t, err)
		assert.True(t, success)
		assert.Equal(t, int32(1), game.Score)
		repo.AssertExpectations(t)
	})

	t.Run("successful tap on type 'b'", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		game := &domain.Game{
			ObjectSeq: []domain.GameObject{
				{
					ID:        "obj1",
					Type:      "b",
					Timestamp: time.Now().Add(-time.Minute),
				},
			},
			Score: 0,
		}

		repo.On("GetGame", "game1").Return(game, nil)
		repo.On("UpdateGame", mock.AnythingOfType("*domain.Game")).Return(nil)

		success, err := service.Tap("game1", "obj1", time.Now())

		assert.NoError(t, err)
		assert.True(t, success)
		assert.Equal(t, int32(-5), game.Score)
		repo.AssertExpectations(t)
	})

	t.Run("game not found", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		repo.On("GetGame", "game1").Return(nil, assert.AnError)

		success, err := service.Tap("game1", "obj1", time.Now())

		assert.Error(t, err)
		assert.False(t, success)
		repo.AssertExpectations(t)
	})

	t.Run("object not found", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		game := &domain.Game{
			ObjectSeq: []domain.GameObject{},
			Score:     0,
		}

		repo.On("GetGame", "game1").Return(game, nil)

		success, err := service.Tap("game1", "nonexistent", time.Now())

		assert.NoError(t, err)
		assert.False(t, success)
		assert.Equal(t, int32(0), game.Score)
		repo.AssertExpectations(t)
	})
}

func TestEndGame(t *testing.T) {
	t.Run("successful game end", func(t *testing.T) {
		t.Skip("Skipping this specific test case")
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		game := &domain.Game{Score: 10}
		repo.On("GetGame", "game1").Return(game, nil)
		repo.On("UpdateGame", mock.AnythingOfType("*domain.Game")).Return(nil)

		score, err := service.EndGame("game1")

		assert.NoError(t, err)
		assert.Equal(t, int32(10), score)
		assert.NotZero(t, game.EndTime)
		repo.AssertExpectations(t)
	})

	t.Run("game not found", func(t *testing.T) {
		t.Skip("Skipping this specific test case")
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		repo.On("GetGame", "game1").Return(nil, assert.AnError)

		score, err := service.EndGame("game1")

		assert.Error(t, err)
		assert.Equal(t, int32(0), score)
		repo.AssertExpectations(t)
	})

	t.Run("update game fails", func(t *testing.T) {
		repo := new(MockGameRepository)
		encrypter := new(MockEncrypter)
		service := NewGameService(repo, encrypter)

		game := &domain.Game{Score: 10}
		repo.On("GetGame", "game1").Return(game, nil)
		repo.On("UpdateGame", mock.AnythingOfType("*domain.Game")).Return(errors.New("update failed"))

		score, err := service.EndGame("game1")

		assert.Error(t, err)
		assert.Equal(t, int32(0), score)
		repo.AssertExpectations(t)
	})
}
