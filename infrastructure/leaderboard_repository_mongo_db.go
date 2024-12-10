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
	opts := options.Client().ApplyURI("mongodb+srv://" + os.Getenv("MONGO_DB_USER") + ":" + os.Getenv("MONGO_DB_PASSWORD") + "@" + os.Getenv("MONGO_DB_URL") + "/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
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
	doc := bson.M{"$set": bson.M{
		"ID":      table.ID,
		"Entries": table.Entries,
	}}
	// Check if user already exists
	filter := bson.M{"ID": table.ID}
	opts := options.Update().SetUpsert(true)
	// Use UpdateOne with upsert to either insert new or update existing
	_, err := repo.collection.UpdateOne(
		context.TODO(),
		filter,
		doc,
		opts,
	)
	if err != nil {
		return fmt.Errorf("failed to save leaderboard: %w", err)
	}
	return nil
}

// GetLeaderboard retrieves a table by its ID
func (repo *MongoDbLeaderboardRepository) GetLeaderboard(tableID string) (*domain.Table, error) {
	fmt.Println("")
	table := &domain.Table{}
	err := repo.collection.FindOne(context.TODO(), bson.M{"ID": tableID}).Decode(&table)
	if err != nil {
		fmt.Println("MongoDbLeaderboardRepository", "GetLeaderboard", "error", tableID, err)
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("table not found")
		}
		return nil, err
	}
	fmt.Println("MongoDbLeaderboardRepository", "GetLeaderboard", table)
	fmt.Println("")
	return table, nil
}

// UpdateLeaderboard updates an existing table in MongoDB
func (repo *MongoDbLeaderboardRepository) UpdateLeaderboard(table *domain.Table) error {
	fmt.Println("")
	fmt.Println("MongoDbLeaderboardRepository", "UpdateLeaderboard", table.ID)
	ctx := context.Background()
	update := bson.M{"$set": bson.M{
		"ID":      table.ID,
		"Entries": table.Entries,
	}}
	result, err := repo.collection.UpdateOne(
		ctx,
		bson.M{"ID": table.ID},
		update,
	)
	if err != nil {
		return fmt.Errorf("failed to update leaderboard: %w", err)
	}
	if result.MatchedCount == 0 {
		return errors.New("table not found")
	}
	fmt.Println("MongoDbLeaderboardRepository", "UpdateLeaderboard", result.ModifiedCount)
	fmt.Println("")
	return nil
}

// AddEntryToLeaderboard adds a new entry to an existing leaderboard
func (repo *MongoDbLeaderboardRepository) AddEntryToLeaderboard(table *domain.Table, entry *domain.GameEntry) error {
	result, err := repo.collection.UpdateOne(
		context.TODO(),
		bson.M{"ID": table.ID},
		bson.M{
			"$push": bson.M{
				"Entries": &entry,
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
