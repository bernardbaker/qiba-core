package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

// type User struct {
// 	UserId       int64
// 	Username     string
// 	FirstName    string
// 	LastName     string
// 	LanguageCode string
// 	IsBot        bool
// 	BonusGames   int64
// }

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserId       int64              `bson:"UserId"`
	BonusGames   int64              `bson:"BonusGames"`
	FirstName    string             `bson:"FirstName"`
	IsBot        bool               `bson:"IsBot"`
	LanguageCode string             `bson:"LanguageCode"`
	Username     string             `bson:"Username"`
	LastName     string             `bson:"lastName"`
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
		BonusGames:   user.BonusGames,
	}
}
