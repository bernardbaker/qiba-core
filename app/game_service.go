package app

import (
	"encoding/json"
	"strconv"
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

func (s *GameService) StartGame(userId string) (string, string, string, error) {
	game := domain.NewGame(userId)
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

func (s *GameService) CanPlay(user domain.User, timestamp time.Time) bool {
	// get all games for user
	games, err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))
	if err != nil {
		return false
	}
	// check if the last game was played more than 24 hours ago
	// if so, return true
	if len(games) == 0 {
		return true
	}
	return timestamp.After(games[len(games)-1].EndTime.Add(30 * time.Second))
}
