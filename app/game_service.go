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
	var possibleNewUser *domain.User

	u, getErr := s.userRepo.Get(userId)
	if getErr != nil {
		fmt.Println("StartGame", getErr)
		possibleNewUser = domain.NewUser(user)
	} else {
		possibleNewUser = u
	}
	saveErr := s.userRepo.Update(possibleNewUser)
	if saveErr != nil {
		fmt.Println("error saving user: ", saveErr)
		return "", "", "", saveErr
	}
	// Encrypt game data and generate HMAC
	// encryptedData, hmac, err := s.encrypter.EncryptGameData(game.ObjectSeq)
	// if err != nil {
	// 	return "", "", "", err
	// }

	// count = 1

	return "", "", game.ID, nil
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
	fmt.Println("EndGame with game ID", game.ID)
	game.EndTime = time.Now().UTC()
	updateError := s.repo.UpdateGame(game)
	if updateError != nil {
		fmt.Println("EndGame", "updateError = s.repo.UpdateGame(game)", updateError)
		game.Score = 0
	}
	// Reset playsLeft counter
	return game.Score, nil
}

func (s *GameService) CanPlay(user domain.User) bool {
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
		}

		if getUser.BonusGames > 0 {
			return true
		}

		// check if the last game was played more than 24 hours ago
		lastGame := games[0]
		fmt.Println("")
		fmt.Println("lastGame")
		fmt.Println(lastGame)
		fmt.Println("")
		timeReference := lastGame.EndTime.UTC().Add(1 * time.Minute)
		fmt.Println("r", timeReference)
		// use server timestamp instead of what is sent
		time := time.Now().UTC()
		latestGame := time.UTC().After(timeReference)
		fmt.Println("c", time.UTC())

		fmt.Println("CanPlay", "return latestGame", latestGame)
		fmt.Println("")
		return latestGame
	}
}

func (s *GameService) AddBonusGame(user domain.User) (bool, error) {
	fmt.Println("AddBonusGame userId", user.UserId)
	user.BonusGames++
	fmt.Println("AddBonusGame user.BonusGames", user.BonusGames)
	saveErr := s.userRepo.Update(&user)
	if saveErr != nil {
		fmt.Println("AddBonusGame saveErr := s.userRepo.Update(&user)", saveErr)
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
		for _, entry := range mocks.GenerateMockData(100) {
			addError := s.leaderboardRepo.AddEntryToLeaderboard(table, entry)
			if addError != nil {
				fmt.Println(addError)
			}
		}
		domain.OrderLeaderboard(table)
	}
}

func (s *GameService) GetLeaderboard(name string, user *domain.User) (string, string, error) {
	table, err := s.leaderboardRepo.GetLeaderboard(name)
	// if table is nil, create a new one
	if table == nil {
		fmt.Println("GameService GetLeaderboard table is nil")
		return "", "", err
	}

	// Sort the scores where duplicate scores are sorted by first to score that amount
	domain.OrderLeaderboard(table)

	// create a map of string to store details
	results := make([]domain.LeaderboardEntry, 0, 100)

	usersScore := make([]domain.LeaderboardEntry, 0, 1)
	didntFindUser := true

	fmt.Println("")
	fmt.Println("GetLeaderboard User", user)
	fmt.Println("")

	// Loop through the entries
	for i, entry := range table.Entries {
		if i > 99 {
			fmt.Println("GetLeaderboard", "i > 99", i)
			break
		}
		results = append(results, domain.LeaderboardEntry{
			Username:  entry.User.Username,
			Score:     entry.Score,
			Timestamp: entry.Timestamp,
		})
		if user != nil && entry.User.UserId == user.UserId {
			fmt.Println("Leaderboard found user with score", entry)
			didntFindUser = false
		}
	}

	if didntFindUser && user != nil {
		fmt.Println("")
		fmt.Println("GetLeaderboard didn't find user && user != nil")
		fmt.Println("")
		for _, entry := range table.Entries {
			if entry.User.UserId == user.UserId {
				usersScore = append(usersScore, domain.LeaderboardEntry{
					Username:  entry.User.Username,
					Score:     entry.Score,
					Timestamp: entry.Timestamp,
				})
				fmt.Println("")
				fmt.Println("if entry.User.UserId == user.UserId")
				fmt.Println(usersScore)
				fmt.Println("")
			}
		}
	}

	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", "", fmt.Errorf("error converting to JSON: %v", err)
	}

	if !didntFindUser && user == nil {
		return string(jsonData), "", nil
	} else {
		userData, err := json.Marshal(usersScore)
		if err != nil {
			return "", "", fmt.Errorf("error converting to JSON: %v", err)
		}
		fmt.Println("Leaderboard found user with score out of top 100 group", userData)
		return string(jsonData), string(userData), nil
	}
}

func (s *GameService) AddToLeaderboard(user domain.User, score int32) (*domain.Table, error) {
	fmt.Println("")
	entry := domain.NewLeaderboardObject(user, score)
	if entry == nil {
		fmt.Println("GameService", "AddToLeaderboard", "entry error", entry)
		return nil, errors.New("entry is nil")
	}
	table, err := s.leaderboardRepo.GetLeaderboard("qiba")
	if err != nil {
		fmt.Println("GameService", "GetLeaderboard", "error", err)
		return nil, err
	}
	addError := s.leaderboardRepo.AddEntryToLeaderboard(table, entry)
	if addError != nil {
		fmt.Println("GameService", "AddEntryToLeaderboard", "addError", addError)
		return nil, addError
	}
	fmt.Println("GameService", "GetLeaderboard", "table", table)
	fmt.Println("")
	return table, nil
}

func (s *GameService) UpdateLeaderboard(table *domain.Table) error {
	err := s.leaderboardRepo.SaveLeaderboard(table)
	if err != nil {
		return err
	}
	return nil
}

func (s *GameService) GetUserScore(leaderboard *domain.Table, userId int64) (*domain.GameEntry, error) {
	for _, entry := range leaderboard.Entries {
		if entry.User.UserId == userId {
			return &entry, nil
		}
	}
	return nil, errors.New("user not found")
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

	count := 1

	// get all games for user
	games, err := s.repo.GetGamesByUser(userId)

	if err != nil {
		fmt.Println("MaxPlays err := s.repo.GetGamesByUser(userId)", err)
	}

	if len(games) > 0 {
		if u.BonusGames > 0 {
			fmt.Println("MaxPlays u.BonusGames > 0", u.BonusGames)
			count += int(u.BonusGames)
		}
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

	// get all games for user
	games, err := s.repo.GetGamesByUser(strconv.FormatInt(u.UserId, 10))

	if err != nil {
		fmt.Println("PlayCount err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))", err)
	}

	count := 0

	dayStart := domain.StartOfDay(time.Now())
	dayEnd := domain.EndOfDay(time.Now())
	for _, game := range games {
		if game.StartTime.UnixMilli() >= dayStart.UnixMilli() && game.EndTime.UnixMilli() <= dayEnd.UnixMilli() {
			count++
		}
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

	if u != nil && u.BonusGames > 0 {
		fmt.Println("PlaysLeft u.BonusGames > 0", u.BonusGames)
		count += int(u.BonusGames)
	}
	// get all games for user
	games, err := s.repo.GetGamesByUser(strconv.FormatInt(u.UserId, 10))

	if err != nil {
		fmt.Println("PlaysLeft err := s.repo.GetGamesByUser(strconv.FormatInt(user.UserId, 10))", err)
	}

	if len(games) == 0 {
		count++
	} else {
		// check if the last game was played more than 24 hours ago
		lastGame := games[0]
		fmt.Println("")
		fmt.Println("lastGame")
		fmt.Println(lastGame)
		fmt.Println("")
		timeReference := lastGame.EndTime.UTC().Add(1 * time.Minute)
		fmt.Println("r", timeReference)
		// use server timestamp instead of what is sent
		time := time.Now().UTC()
		latestGame := time.UTC().After(timeReference)
		fmt.Println("c", time.UTC())
		if latestGame {
			count++
		}
	}

	count--
	if count < 0 {
		count = 0
	}

	return int32(count)
}
