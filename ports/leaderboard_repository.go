package ports

import "github.com/bernardbaker/qiba.core/domain"

// UserRepository defines the repository interface for user data
type LeaderboardRepository interface {
	SaveLeaderboard(leaderboard *domain.Table) error
	GetLeaderboard(name string) (*domain.Table, error)
	AddEntryToLeaderboard(leaderboard *domain.Table, entry *domain.GameEntry) error
	UpdateLeaderboard(table *domain.Table) error
}
