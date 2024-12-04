package app

import (
	"errors"
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
	u := s.repo.Get(owner)
	if u == nil {
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

func (s *ReferralService) Get(objectID string) (*domain.Referral, bool) {
	obj := s.repo.Get(objectID)
	if obj == nil {
		fmt.Println(errors.New("service get referral not found"))
		return nil, false
	}
	return obj, true
}

func (s *ReferralService) Update(from domain.User, to domain.User, gameService GameService) (bool, bool) {
	owner := strconv.FormatInt(from.UserId, 10)
	obj := s.repo.Get(owner)
	if obj == nil {
		fmt.Println(errors.New("service update referral not found"))
		return false, false
	}
	// check the referrals
	hasReferral := false
	for _, referral := range obj.Referrals {
		// if to has already had a referral from from, then skip
		if referral.To.UserId == to.UserId && referral.From.UserId == from.UserId {
			fmt.Println("To has already had a referral from From...")
			hasReferral = true
		}
	}
	// if so, return error
	if !hasReferral {
		// otherwise, add the new referral
		obj.Referrals = append(obj.Referrals, *domain.NewReferralObject(from, to))
		// store the referral
		updateError := s.repo.Update(obj)
		if !updateError {
			fmt.Println(errors.New("service update referral not found"))
			return false, false
		}

		fmt.Println("")
		fmt.Println("success, addBonusErr := s.gameService.AddBonusGame(from)", from)
		success, addBonusErr := gameService.AddBonusGame(from)
		if addBonusErr != nil {
			fmt.Println("addBonusErr", addBonusErr)
			return success, true
		}
	}

	bonusGames := from.BonusGames
	bonusGames++
	from.BonusGames = bonusGames

	gameService.userRepo.Update(&from)

	return false, false
}
