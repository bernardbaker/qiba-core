package ports

import "github.com/bernardbaker/qiba.core/domain"

// GameRepository defines the repository interface for game data
type ReferralRepository interface {
	Save(object *domain.Referral) error
	Get(objectID string) (*domain.Referral, error)
	Update(object *domain.Referral) error
}
