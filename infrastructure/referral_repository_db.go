package infrastructure

import (
	"errors"
	"fmt"
	"sync"

	"github.com/bernardbaker/qiba.core/domain"
)

type InMemoryReferralRepository struct {
	store map[string]*domain.Referral
	mutex sync.RWMutex
}

func NewInMemoryReferralRepository() *InMemoryReferralRepository {
	return &InMemoryReferralRepository{
		store: make(map[string]*domain.Referral),
	}
}

// Save stores a new referral in the in-memory map
func (repo *InMemoryReferralRepository) Save(obj *domain.Referral) error {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.store[obj.ID] = obj
	return nil
}

// Get retrieves a referral by its ID
func (repo *InMemoryReferralRepository) Get(objId string) *domain.Referral {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	obj, exists := repo.store[objId]
	if !exists {
		fmt.Println(errors.New("referral not found"))
		return nil
	}
	return obj
}

// Update updates an existing referral in the in-memory map
func (repo *InMemoryReferralRepository) Update(obj *domain.Referral) bool {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	_, exists := repo.store[obj.ID]
	if !exists {
		fmt.Println(errors.New("referral not found"))
		return false
	}
	repo.store[obj.ID] = obj
	return true
}
