package app

import (
	"encoding/json"
	"time"

	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/ports"
)

type GameService struct {
	repo      ports.GameRepository
	encrypter ports.Encrypter
}

func NewGameService(repo ports.GameRepository, encrypter ports.Encrypter) *GameService {
	return &GameService{repo: repo, encrypter: encrypter}
}

func (s *GameService) StartGame() (string, string, string, error) {
	game := domain.NewGame()
	err := s.repo.SaveGame(game)
	if err != nil {
		return "", "", "", err
	}

	// Encrypt game data and generate HMAC
	encryptedData, hmac, err := s.encrypter.EncryptGameData(game.ObjectSeq)
	if err != nil {
		return "", "", "", err
	}
	return encryptedData, hmac, game.ID, nil
}

func (s *GameService) Spawn(gameID string) (string, error) {
	game, err := s.repo.GetGame(gameID)
	if err != nil {
		return "", err
	}

	game.GenerateObjectSequence()

	s.repo.SaveGame(game)

	json, err := json.Marshal(game.ObjectSeq[len(game.ObjectSeq)-1])

	if err != nil {
		return "", err
	}

	return string(json), nil
}

func (s *GameService) Tap(gameID, objectID string, timestamp time.Time) (bool, error) {
	game, err := s.repo.GetGame(gameID)
	if err != nil {
		return false, err
	}

	// Verify object ID and timestamp
	for _, obj := range game.ObjectSeq {
		if obj.ID == objectID && time.Now().After(obj.Timestamp) {
			if obj.Type == "a" {
				game.Score++
			} else {
				game.Score = game.Score - 5
			}
			return true, s.repo.UpdateGame(game)
		}
	}
	return false, nil
}

func (s *GameService) EndGame(gameID string) (int32, error) {
	game, err := s.repo.GetGame(gameID)
	if err != nil {
		return 0, err
	}
	updateError := s.repo.UpdateGame(game)
	if updateError != nil {
		game.Score = 0
	}
	game.EndTime = time.Now()
	return game.Score, updateError
}
