package domain

type GameRepository interface {
	SaveGame(game *Game) error
	GetGame(gameID string) (*Game, error)
	UpdateGame(game *Game) error
	GenerateObjectSequence(gameID string) (*Game, error)
}
