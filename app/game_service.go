package app

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/ports"
)

type GameService struct {
	repo      ports.GameRepository
	userRepo  ports.UserRepository
	encrypter ports.Encrypter
}

func NewGameService(repo ports.GameRepository, userRepo ports.UserRepository, encrypter ports.Encrypter) *GameService {
	return &GameService{repo: repo, userRepo: userRepo, encrypter: encrypter}
}

func (s *GameService) StartGame(userId string, user domain.User) (string, string, string, error) {
	game := domain.NewGame(userId)
	err := s.repo.SaveGame(game)
	if err != nil {
		return "", "", "", err
	}
	possibleNewUser := domain.NewUser(user)
	saveErr := s.userRepo.Save(possibleNewUser)
	if saveErr != nil {
		fmt.Println("Error saving user: ", saveErr)
		return "", "", "", saveErr
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
	if len(games) == 0 {
		return true
	}
	// check if the last game was played more than 24 hours ago
	latestGame := timestamp.After(games[len(games)-1].EndTime.Add(30 * time.Second))

	if !latestGame {

		// debugging
		fmt.Println("")
		getUser, err := s.userRepo.Get(strconv.FormatInt(user.UserId, 10))
		if err != nil {
			fmt.Println("CanPlay", "err = s.userRepo.Get(strconv.FormatInt(user.UserId, 10))", err)
			fmt.Println(err)
		}

		fmt.Println("CanPlay", "!latestGame", latestGame)
		if getUser.BonusGames > 0 {
			fmt.Println("CanPlay", "getUser..BonusGames > 0", getUser.BonusGames)

			getUser.BonusGames--
			fmt.Println("CanPlay", "getUser..BonusGames is now", getUser.BonusGames)
			err = s.userRepo.Save(getUser)
			if err != nil {
				fmt.Println("CanPlay", "err = s.userRepo.Save(getUser)", err)
				fmt.Println(err)
			}

			latestGame = true
			fmt.Println("CanPlay", "latestGame = true", latestGame)
		}
	}
	fmt.Println("CanPlay", "return latestGame", latestGame)
	fmt.Println("")
	return latestGame
}

func (s *GameService) AddBonusGame(user domain.User) (bool, error) {
	// convert user.UserId to a string
	userId := strconv.FormatInt(user.UserId, 10)
	// debugging
	fmt.Println("AddBonusGame userId", userId)
	u, err := s.userRepo.Get(userId)
	if err != nil {
		fmt.Println("AddBonusGame err := s.userRepo.Get(userId)", err)
		return false, err
	}
	fmt.Println("AddBonusGame u", u)
	u.BonusGames++
	fmt.Println("AddBonusGame u.BonusGames", u.BonusGames)
	saveErr := s.userRepo.Update(u)
	if saveErr != nil {
		fmt.Println("AddBonusGame saveErr := s.userRepo.Save(u)", saveErr)
		return false, saveErr
	}
	return true, nil
}

func (s *GameService) AddUser(user domain.User) (bool, error) {
	possibleNewUser := domain.NewUser(user)
	err := s.userRepo.Save(possibleNewUser)
	if err != nil {
		fmt.Println("Error saving possible new user: ", err)
		return false, err
	}
	return true, nil
}

func (s *GameService) GetBonusGames(user domain.User) (string, bool) {
	// convert user.UserId to a string
	userId := strconv.FormatInt(user.UserId, 10)
	// debugging
	fmt.Println("GetBonusGames userId", userId)
	u, err := s.userRepo.Get(userId)
	if err != nil {
		fmt.Println("GetBonusGames u, err := s.userRepo.Get(userId)", err)
		return string(0), false
	}
	count := strconv.FormatInt(u.BonusGames, 10)
	fmt.Println("count", count)
	return count, true
}
