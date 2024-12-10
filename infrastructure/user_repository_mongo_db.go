package infrastructure

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/bernardbaker/qiba.core/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// In infrastructure/mongodb_repository.go

type MongoDbUserRepository struct {
	client *mongo.Client
	// add any other fields you need, like collection name
	collection *mongo.Collection
}

// Constructor
func NewMongoDbUserRepository() *MongoDbUserRepository {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGO_DB_USER") + ":" + os.Getenv("MONGO_DB_PASSWORD") + "@" + os.Getenv("MONGO_DB_URL") + "/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println("User repository - connection to MongoDB failed!")
	}
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("User repository - Pinged your deployment. You successfully connected to MongoDB!")

	// Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &MongoDbUserRepository{
		client:     client,
		collection: client.Database("qiba-game").Collection("users"),
	}
}

// Implement all methods required by the ports.UserRepository interface
func (r *MongoDbUserRepository) Get(id string) (*domain.User, error) {
	var user domain.User
	userId, _ := strconv.ParseInt(id, 10, 64)
	fmt.Println("")
	fmt.Println("User ", userId)
	fmt.Println("")
	err := r.collection.FindOne(context.TODO(), bson.M{"UserId": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoDbUserRepository) Save(user *domain.User) error {
	// Convert the domain User to a BSON-friendly format
	userDoc := bson.M{"$set": bson.M{
		"UserId":       user.UserId,
		"BonusGames":   user.BonusGames,
		"FirstName":    user.FirstName,
		"IsBot":        user.IsBot,
		"LanguageCode": user.LanguageCode,
		"lastName":     user.LastName,
		"Username":     user.Username,
	}}

	// Check if user already exists
	filter := bson.M{"UserId": user.UserId}
	opts := options.Update().SetUpsert(true)

	// Use UpdateOne with upsert to either insert new or update existing
	_, err := r.collection.UpdateOne(
		context.TODO(),
		filter,
		userDoc,
		opts,
	)

	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	return nil
}

func (r *MongoDbUserRepository) Update(user *domain.User) error {
	_, err := r.collection.UpdateOne(
		context.TODO(),
		bson.M{"UserId": user.UserId},
		bson.M{"$set": bson.M{
			"UserId":       user.UserId,
			"BonusGames":   user.BonusGames,
			"FirstName":    user.FirstName,
			"IsBot":        user.IsBot,
			"LanguageCode": user.LanguageCode,
			"lastName":     user.LastName,
			"Username":     user.Username,
		}},
	)
	return err
}
