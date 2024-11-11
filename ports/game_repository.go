package ports

import "github.com/bernardbaker/qiba.core/domain"

// GameRepository defines the repository interface for game data
type GameRepository interface {
	SaveGame(game *domain.Game) error
	GetGame(gameID string) (*domain.Game, error)
	UpdateGame(game *domain.Game) error
}
