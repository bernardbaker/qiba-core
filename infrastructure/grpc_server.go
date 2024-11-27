package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/bernardbaker/qiba.core/app"
	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/proto"
)

type GameServer struct {
	proto.UnimplementedGameServiceServer
	service *app.GameService
}

func NewGameServer(service *app.GameService) *GameServer {
	return &GameServer{service: service}
}

func (s *GameServer) StartGame(ctx context.Context, req *proto.StartGameRequest) (*proto.StartGameResponse, error) {
	id := strconv.FormatInt(req.User.UserId, 10)
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	encryptedData, hmac, gameID, err := s.service.StartGame(id, user)
	if err != nil {
		return nil, err
	}
	return &proto.StartGameResponse{EncryptedGameData: encryptedData, Hmac: hmac, GameId: gameID}, nil
}

func (s *GameServer) Spawn(ctx context.Context, req *proto.SpawnRequest) (*proto.SpawnResponse, error) {
	data, err := s.service.Spawn(req.GameId)
	if err != nil {
		return nil, err
	}
	return &proto.SpawnResponse{Data: data}, nil
}

func (s *GameServer) Tap(ctx context.Context, req *proto.TapRequest) (*proto.TapResponse, error) {
	timestamp, _ := time.Parse(time.RFC3339, req.Timestamp)
	success, err := s.service.Tap(req.GameId, req.ObjectId, timestamp)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &proto.TapResponse{Success: success}, nil
}

func (s *GameServer) EndGame(ctx context.Context, req *proto.EndGameRequest) (*proto.EndGameResponse, error) {
	score, err := s.service.EndGame(req.GameId)
	if err != nil {
		return nil, err
	}
	leaderboard, err := s.service.GetLeaderboard("qiba")
	fmt.Println("gRPC server EndGame", err)
	fmt.Println("leaderboards", leaderboard)
	fmt.Println("req.User", req.User)

	// create a user
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	fmt.Println("EndGame user", user)
	table, err := s.service.AddToLeaderboard(user, score)
	if err != nil {
		return nil, errors.New("failed to add user to leaderboard")
	}
	fmt.Println("after leaderboard table", len(table.Entries))
	saveErr := s.service.SaveLeaderboard(table)
	if saveErr != nil {
		return nil, errors.New("failed to save leaderboard")
	}
	// TODO Using the returned user update the Table
	return &proto.EndGameResponse{Score: score}, nil
}

func (s *GameServer) CanPlay(ctx context.Context, req *proto.CanPlayGameRequest) (*proto.CanPlayGameResponse, error) {
	// Convert string to int64 first if req.Timestamp is a string
	milliseconds, err := strconv.ParseInt(req.Timestamp, 10, 64)
	if err != nil {
		fmt.Println(err)
		return &proto.CanPlayGameResponse{Success: false}, nil
	}
	user := domain.User{
		UserId: req.User.UserId,
	}
	timestamp := time.UnixMilli(milliseconds).UTC()
	success := s.service.CanPlay(user, timestamp)
	return &proto.CanPlayGameResponse{Success: success}, nil
}

// TODO: finish this off
func (s *GameServer) Leaderboard(ctx context.Context, req *proto.LeaderboardRequest) (*proto.LeaderboardResponse, error) {
	fmt.Println("")
	fmt.Println("Leaderboard")

	// create a user
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	fmt.Println("req.User", user)
	// Get the domain table
	jsonString, err := s.service.GetLeaderboard("qiba")
	if err != nil {
		return nil, err
	}
	return &proto.LeaderboardResponse{Success: true, Table: jsonString}, nil
}

func (s *GameServer) GameTime(ctx context.Context, req *proto.GameTimeRequest) (*proto.GameTimeResponse, error) {
	value := s.service.GameTime()
	return &proto.GameTimeResponse{Success: true, Time: value}, nil
}

func (s *GameServer) MaxPlays(ctx context.Context, req *proto.MaxPlaysRequest) (*proto.MaxPlaysResponse, error) {
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	value := s.service.MaxPlays(user)
	return &proto.MaxPlaysResponse{Success: true, Value: value}, nil
}

func (s *GameServer) PlayCount(ctx context.Context, req *proto.PlayCountRequest) (*proto.PlayCountResponse, error) {
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	value := s.service.PlayCount(user)
	return &proto.PlayCountResponse{Success: true, Value: value}, nil
}

func (s *GameServer) PlaysLeft(ctx context.Context, req *proto.PlaysLeftRequest) (*proto.PlaysLeftResponse, error) {
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	value := s.service.PlaysLeft(user)
	return &proto.PlaysLeftResponse{Success: true, Value: value}, nil
}

type ReferralServer struct {
	proto.UnimplementedReferralServiceServer
	service     *app.ReferralService
	gameService *app.GameService
}

func NewReferralServer(service *app.ReferralService, gameService *app.GameService) *ReferralServer {
	return &ReferralServer{service: service, gameService: gameService}
}

func (s *ReferralServer) Referral(ctx context.Context, req *proto.ReferralRequest) (*proto.ReferralResponse, error) {
	fmt.Println("req.User", req.User)
	// Create a new user
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	_, addErr := s.gameService.AddUser(user)
	if addErr != nil {
		return nil, addErr
	}
	// Add the user to the service
	createErr := s.service.Create(req.User.UserId)
	if createErr != nil {
		return nil, createErr
	}
	// debugging
	fmt.Println(s.service.Get(strconv.FormatInt(req.User.UserId, 10)))
	//
	return &proto.ReferralResponse{Success: true}, nil
}

func (s *ReferralServer) AcceptReferral(ctx context.Context, req *proto.AcceptReferralRequest) (*proto.AcceptReferralResponse, error) {
	// debugging
	fmt.Println("AcceptReferral")
	fmt.Println(s.service.Get(strconv.FormatInt(req.From.UserId, 10)))
	from := domain.User{
		UserId:       req.From.UserId,
		Username:     req.From.Username,
		FirstName:    req.From.FirstName,
		LastName:     req.From.LastName,
		LanguageCode: req.From.LanguageCode,
		IsBot:        req.From.IsBot,
	}
	to := domain.User{
		UserId:       req.To.UserId,
		Username:     req.To.Username,
		FirstName:    req.To.FirstName,
		LastName:     req.To.LastName,
		LanguageCode: req.To.LanguageCode,
		IsBot:        req.To.IsBot,
	}
	success, err := s.service.Update(from, to, *s.gameService)
	if !err {
		return nil, errors.New("gRPC server accept referral update error")
	}
	return &proto.AcceptReferralResponse{Success: success}, nil
}

func (s *ReferralServer) ReferralStatistics(ctx context.Context, req *proto.ReferralStatisticsRequest) (*proto.ReferralStatisticsResponse, error) {
	objects, err := s.service.Get(strconv.FormatInt(req.User.UserId, 10))
	if !err {
		fmt.Println(errors.New("gRPC server referral statistics error"))
		return nil, nil
	}
	count := int64(len(objects.Referrals))
	user := domain.User{
		UserId:       req.User.UserId,
		Username:     req.User.Username,
		FirstName:    req.User.FirstName,
		LastName:     req.User.LastName,
		LanguageCode: req.User.LanguageCode,
		IsBot:        req.User.IsBot,
	}
	bCount, bErr := s.gameService.GetBonusGames(user)
	if !bErr {
		fmt.Println(errors.New("gRPC server referral statistics get bonus games error"))
		return nil, nil
	}
	return &proto.ReferralStatisticsResponse{Success: true, Count: count, BonusCount: bCount}, nil
}
