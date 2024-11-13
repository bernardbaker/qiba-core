package app

import (
	"fmt"
	"strconv"

	"github.com/bernardbaker/qiba.core/domain"
	"github.com/bernardbaker/qiba.core/ports"
)

type ReferralService struct {
	repo ports.ReferralRepository
}

func NewReferralService(repo ports.ReferralRepository) *ReferralService {
	return &ReferralService{repo: repo}
}

func (s *ReferralService) Create(user int64) error {
	owner := strconv.FormatInt(user, 10)
	fmt.Println("owner", owner)
	_, err := s.repo.Get(owner)
	if err != nil {
		saveErr := s.repo.Save(domain.NewReferral(owner))
		if saveErr != nil {
			return saveErr
		}
	}
	return nil
}

func (s *ReferralService) Save(object *domain.Referral) error {
	err := s.repo.Save(object)
	if err != nil {
		return err
	}
	return nil
}

func (s *ReferralService) Get(objectID string) (*domain.Referral, error) {
	obj, err := s.repo.Get(objectID)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *ReferralService) Update(from domain.User, to domain.User) error {
	owner := strconv.FormatInt(from.UserId, 10)
	obj, err := s.repo.Get(owner)
	if err != nil {
		return err
	}
	// check the referrals
	for _, referral := range obj.Referrals {
		// if to has already had a referral from from, then skip
		if referral.To.UserId == to.UserId && referral.From.UserId == from.UserId {
			fmt.Println("To has already had a referral from From...")
			return nil
		}
	}
	// otherwise, add the new referral
	obj.Referrals = append(obj.Referrals, *domain.NewReferralObject(from, to))
	// store the referral
	updateError := s.repo.Update(obj)
	if updateError != nil {
		return updateError
	}
	return nil
}
