package domain

type User struct {
	UserId       int64
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
	IsBot        bool
	BonusGames   int64
}

// Generate a new game with random object sequence
func NewUser(user User) *User {
	return &User{
		UserId:       user.UserId,
		Username:     user.Username,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		LanguageCode: user.LanguageCode,
		IsBot:        user.IsBot,
	}
}
