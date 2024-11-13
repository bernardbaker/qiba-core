package domain

type ReferralRepository interface {
	Save(object *Referral) error
	Get(objectID string) (*Referral, error)
	Update(object *Referral) error
}
