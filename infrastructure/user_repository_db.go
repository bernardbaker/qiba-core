package infrastructure

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/bernardbaker/qiba.core/domain"
)

type InMemoryUserRepository struct {
	users map[string]*domain.User
	mutex sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

// Save stores a new user in the in-memory map
func (repo *InMemoryUserRepository) Save(user *domain.User) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	userIdStr := strconv.FormatInt(user.UserId, 10)
	_, exists := repo.users[userIdStr]
	if !exists {
		repo.users[userIdStr] = user
	}
	return nil
}

// Get retrieves a user by its ID
func (repo *InMemoryUserRepository) Get(userID string) (*domain.User, error) {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	fmt.Println("NewInMemoryUserRepository Get", repo.users)
	user, exists := repo.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// Update updates an existing user in the in-memory map
func (repo *InMemoryUserRepository) Update(user *domain.User) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	userIdStr := strconv.FormatInt(user.UserId, 10)
	_, exists := repo.users[userIdStr]
	if !exists {
		return errors.New("User Repository Update user not found")
	}
	repo.users[userIdStr] = user
	return nil
}
