package ports

import (
	"github.com/bernardbaker/qiba.core/domain"
)

// UserRepository defines the repository interface for user data
type UserRepository interface {
	Save(obj *domain.User) error
	Get(objID string) (*domain.User, error)
	Update(obj *domain.User) error
}
