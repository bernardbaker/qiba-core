package domain

type User struct {
	UserId       int64
	Username     string
	FirstName    string
	LastName     string
	LanguageCode string
	IsBot        bool
}
