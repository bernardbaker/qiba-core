package domain

import (
	"time"

	"github.com/google/uuid"
)

type Referral struct {
	ID         string `bson:"id"`
	Referrals  []ReferralObject
	CreateTime time.Time
}

type ReferralObject struct {
	ID         string
	From       User
	To         User
	AcceptTime time.Time
	Expired    bool
}

// Generate a new game with random object sequence
func NewReferral(owner string) *Referral {
	referral := &Referral{
		ID:         owner,
		Referrals:  []ReferralObject{},
		CreateTime: time.Now(),
	}
	return referral
}

// Generates a new referral
func NewReferralObject(from User, to User) *ReferralObject {
	referralObject := &ReferralObject{
		ID:         uuid.New().String(),
		From:       from,
		To:         to,
		AcceptTime: time.Time{},
		Expired:    true,
	}
	return referralObject
}
