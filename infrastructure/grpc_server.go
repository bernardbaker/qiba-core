package infrastructure

import (
	"context"
	"time"

	"github.com/bernardbaker/qiba.core/app"
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
	encryptedData, hmac, gameID, err := s.service.StartGame()
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
		return nil, err
	}
	return &proto.TapResponse{Success: success}, nil
}

func (s *GameServer) EndGame(ctx context.Context, req *proto.EndGameRequest) (*proto.EndGameResponse, error) {
	score, err := s.service.EndGame(req.GameId)
	if err != nil {
		return nil, err
	}
	return &proto.EndGameResponse{Score: score}, nil
}
