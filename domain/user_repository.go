package domain

type UserRepository interface {
	Save(object *User) error
	Get(objectID string) (*User, error)
	Update(object *User) error
}
