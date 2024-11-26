package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/mocks"
	"github.com/bernardbaker/qiba.core/ports"
)

type GameService struct {
	repo            ports.GameRepository
	userRepo        ports.UserRepository
	leaderboardRepo ports.LeaderboardRepository
	encrypter       ports.Encrypter
}

func NewGameService(repo ports.GameRepository, userRepo ports.UserRepository, leaderboardRepo ports.LeaderboardRepository, encrypter ports.Encrypter) *GameService {
	return &GameService{repo: repo, userRepo: userRepo, leaderboardRepo: leaderboardRepo, encrypter: encrypter}
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

// TODO Return the game object instead
func (s *GameService) EndGame(gameID string) (int32, error) {
	game, err := s.repo.GetGame(gameID)
	if err != nil {
		return 0, err
	}
	game.EndTime = time.Now()
	updateError := s.repo.UpdateGame(game)
	if updateError != nil {
		game.Score = 0
	}
	return game.Score, updateError
}

func (s *GameService) CanPlay(user domain.User, timestamp time.Time) bool {
	// get all games for user
	games, err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))
	if err != nil {
		return false
	}
	if len(games) == 0 {
		fmt.Println("")
		fmt.Println("CanPlay", "len(games) == 0", len(games))
		fmt.Println("")
		return true
	} else {
		// debugging
		fmt.Println("")
		getUser, err := s.userRepo.Get(strconv.FormatInt(user.UserId, 10))
		if err != nil {
			fmt.Println("CanPlay", "err = s.userRepo.Get(strconv.FormatInt(user.UserId, 10))", err)
			fmt.Println(err)
		}

		if getUser.BonusGames > 0 {
			fmt.Println("CanPlay", "getUser..BonusGames > 0", getUser.BonusGames)
			getUser.BonusGames--
			fmt.Println("CanPlay", "getUser..BonusGames is now", getUser.BonusGames)
			err = s.userRepo.Save(getUser)
			if err != nil {
				fmt.Println("CanPlay", "err = s.userRepo.Save(getUser)", err)
				fmt.Println(err)
			}
			return true
		}
		// check if the last game was played more than 24 hours ago
		timeReference := games[len(games)-1].EndTime.UTC().Add(30 * time.Second)
		fmt.Println("r", timeReference)
		latestGame := timestamp.UTC().After(timeReference)
		fmt.Println("c", timestamp.UTC())

		fmt.Println("CanPlay", "return latestGame", latestGame)
		fmt.Println("")
		return latestGame
	}
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

func (s *GameService) CreateLeaderboard(name string, prepopulate bool) {
	leaderboard := domain.NewLeaderboard(name)
	s.leaderboardRepo.SaveLeaderboard(leaderboard)

	if prepopulate {
		table, err := s.leaderboardRepo.GetLeaderboard("qiba")
		if err != nil {
			fmt.Println(err)
		}
		for _, entry := range mocks.GenerateMockData(10000) {
			addError := s.leaderboardRepo.AddEntryToLeaderboard(table, entry)
			if addError != nil {
				fmt.Println(addError)
			}
		}
		domain.OrderLeaderboard(table)
	}
}

func (s *GameService) GetLeaderboard(name string) (string, error) {
	table, err := s.leaderboardRepo.GetLeaderboard(name)
	// if table is nil, create a new one
	if table == nil {
		return "", err
	}
	fmt.Println("")
	fmt.Println("len(table.Entries)", len(table.Entries))
	fmt.Println("")
	// create a map of string to store details
	results := make([]domain.LeaderboardEntry, 0, len(table.Entries))

	// Loop through the entries
	for _, entry := range table.Entries {
		results = append(results, domain.LeaderboardEntry{
			Username:  entry.User.Username,
			Score:     entry.Score,
			Timestamp: entry.Timestamp,
		})
	}
	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("error converting to JSON: %v", err)
	}

	return string(jsonData), nil
}

func (s *GameService) AddToLeaderboard(user domain.User, score int32) (*domain.Table, error) {
	entry := domain.NewLeaderboardObject(user, score)
	if entry == nil {
		return nil, errors.New("entry is nil")
	}
	table, err := s.leaderboardRepo.GetLeaderboard("qiba")
	if err != nil {
		return nil, err
	}
	addError := s.leaderboardRepo.AddEntryToLeaderboard(table, entry)
	if addError != nil {
		return nil, addError
	}
	return table, nil
}

func (s *GameService) SaveLeaderboard(table *domain.Table) error {
	err := s.leaderboardRepo.SaveLeaderboard(table)
	if err != nil {
		return err
	}
	return nil
}

func (s *GameService) GameTime() int32 {
	return 10
}

func (s *GameService) MaxPlays(user domain.User) int32 {
	// convert user.UserId to a string
	userId := strconv.FormatInt(user.UserId, 10)
	// debugging
	fmt.Println("MaxPlays userId", userId)
	u, err := s.userRepo.Get(userId)
	if err != nil {
		fmt.Println("MaxPlays u, err := s.userRepo.Get(userId)", err)
		return 1
	}

	count := 0
	// get all games for user
	games, err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))

	if err != nil {
		fmt.Println("MaxPlays err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))", err)
	}

	if len(games) == 0 {
		count++
	} else {
		// check if the last game was played more than 24 hours ago
		timestamp := time.Now()
		timeReference := games[len(games)-1].EndTime.UTC().Add(30 * time.Second)
		fmt.Println("r", timeReference)
		latestGame := timestamp.UTC().After(timeReference)
		fmt.Println("c", timestamp.UTC())

		if latestGame {
			count++
		}
	}

	if u.BonusGames > 0 {
		fmt.Println("MaxPlays u.BonusGames > 0", u.BonusGames)
		count += int(u.BonusGames)
	}

	fmt.Println("MaxPlays count", count)

	return int32(count)
}

func (s *GameService) PlayCount(user domain.User) int32 {
	// convert user.UserId to a string
	userId := strconv.FormatInt(user.UserId, 10)
	// debugging
	fmt.Println("PlayCount userId", userId)
	u, err := s.userRepo.Get(userId)
	if err != nil {
		fmt.Println("PlayCount u, err := s.userRepo.Get(userId)", err)
		s.userRepo.Save(domain.NewUser(user))
		return 1
	}

	count := 1

	if u.BonusGames > 0 {
		fmt.Println("PlayCount u.BonusGames > 0", u.BonusGames)
		count += int(u.BonusGames)
	}

	return int32(count)
}

func (s *GameService) PlaysLeft(user domain.User) int32 {
	// convert user.UserId to a string
	userId := strconv.FormatInt(user.UserId, 10)
	// debugging
	fmt.Println("")
	fmt.Println("PlaysLeft userId", userId)
	u, err := s.userRepo.Get(userId)
	if err != nil {
		fmt.Println("PlaysLeft u, err := s.userRepo.Get(userId)", err)
		s.userRepo.Save(domain.NewUser(user))
	}

	count := 0

	if u.BonusGames > 0 {
		fmt.Println("PlaysLeft u.BonusGames > 0", u.BonusGames)
		count += int(u.BonusGames)
	}

	// get all games for user
	games, err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))

	if err != nil {
		fmt.Println("PlaysLeft err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))", err)
	}

	if len(games) == 0 {
		count++
	} else {
		// check if the last game was played more than 24 hours ago
		timestamp := time.Now()
		timeReference := games[len(games)-1].EndTime.UTC().Add(30 * time.Second)
		latestGame := timestamp.UTC().After(timeReference)
		if latestGame {
			count++
		}
		count = count - 1
	}

	return int32(count)
}
