package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/bernardbaker/qiba.core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbLeaderboardRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbLeaderboardRepository() *MongoDbLeaderboardRepository {
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	var opts *options.ClientOptions
	if os.Getenv("ENV") == "development" {
		opts = options.Client().ApplyURI("mongodb+srv://bernard:kCxUAkcAKwhFcEyl@qiba-game.umja7yd.mongodb.net/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
	} else {
		opts = options.Client().ApplyURI("mongodb+srv://bernard:kCxUAkcAKwhFcEyl@qiba-game.umja7yd.mongodb.net/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
		// opts = options.Client().ApplyURI("mongodb+srv://bernard:kCxUAkcAKwhFcEyl@qiba-game-pl-0.lbdk6.mongodb.net/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
	}
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		fmt.Println("Leaderboard repository - connection to MongoDB failed!")
	}
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Leaderboard repository - Pinged your deployment. You successfully connected to MongoDB!")

	return &MongoDbLeaderboardRepository{
		client:     client,
		collection: client.Database("qiba-game").Collection("leaderboard"),
	}
}

// SaveLeaderboard stores a new leaderboard in MongoDB
func (repo *MongoDbLeaderboardRepository) SaveLeaderboard(table *domain.Table) error {
	_, err := repo.collection.InsertOne(context.TODO(), table)
	if err != nil {
		return err
	}
	return nil
}

// GetLeaderboard retrieves a table by its ID
func (repo *MongoDbLeaderboardRepository) GetLeaderboard(tableID string) (*domain.Table, error) {
	var table domain.Table
	err := repo.collection.FindOne(context.TODO(), bson.M{"id": tableID}).Decode(&table)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("table not found")
		}
		return nil, err
	}
	return &table, nil
}

// UpdateLeaderboard updates an existing table in MongoDB
func (repo *MongoDbLeaderboardRepository) UpdateLeaderboard(table *domain.Table) error {
	result, err := repo.collection.ReplaceOne(
		context.TODO(),
		bson.M{"id": table.ID},
		table,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("table not found")
	}
	return nil
}

// AddEntryToLeaderboard adds a new entry to an existing leaderboard
func (repo *MongoDbLeaderboardRepository) AddEntryToLeaderboard(table *domain.Table, entry *domain.GameEntry) error {
	result, err := repo.collection.UpdateOne(
		context.TODO(),
		bson.M{"id": table.ID},
		bson.M{
			"$push": bson.M{
				"entries": entry,
			},
		},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("table not found")
	}
	return nil
}
