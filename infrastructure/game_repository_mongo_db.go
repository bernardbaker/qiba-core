package infrastructure

import (
	"cmp"
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/bernardbaker/qiba.core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbGameRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbGameRepository() *MongoDbGameRepository {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	var opts *options.ClientOptions
	if os.Getenv("ENV") == "development" {
		opts = options.Client().ApplyURI("mongodb+srv://bernard:kCxUAkcAKwhFcEyl@qiba-game.umja7yd.mongodb.net/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
	} else {
		opts = options.Client().ApplyURI("mongodb+srv://bernard:kCxUAkcAKwhFcEyl@qiba-game.umja7yd.mongodb.net/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
		// opts = options.Client().ApplyURI("mongodb+srv://bernard:kCxUAkcAKwhFcEyl@qiba-game-pl-0.lbdk6.mongodb.net/?retryWrites=true&w=majority&appName=qiba-game").SetServerAPIOptions(serverAPI)
	}

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		fmt.Println("Game repository - connection to MongoDB failed!")
	}
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Game repository - Pinged your deployment. You successfully connected to MongoDB!")

	return &MongoDbGameRepository{
		client:     client,
		collection: client.Database("qiba-game").Collection("games"),
	}
}

// SaveGame stores a new game in MongoDB
func (repo *MongoDbGameRepository) SaveGame(game *domain.Game) error {
	ctx := context.Background()
	filter := bson.M{"_id": game.ID}

	update := bson.M{"$set": bson.M{
		"EndTime":   game.EndTime,
		"ID":        game.ID,
		"ObjectSeq": game.ObjectSeq,
		"Score":     game.Score,
		"StartTime": game.StartTime,
		"UserID":    game.UserID,
	}}
	opts := options.Update().SetUpsert(true)

	_, err := repo.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save game: %w", err)
	}

	return nil
}

// GetGame retrieves a game by its ID from MongoDB
func (repo *MongoDbGameRepository) GetGame(gameID string) (*domain.Game, error) {
	ctx := context.Background()
	var game domain.Game
	filter := bson.M{"_id": gameID}

	err := repo.collection.FindOne(ctx, filter).Decode(&game)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("game not found with id: %s", gameID)
		}
		return nil, fmt.Errorf("error fetching game: %w", err)
	}

	return &game, nil
}

// UpdateGame updates an existing game in MongoDB
func (repo *MongoDbGameRepository) UpdateGame(game *domain.Game) error {
	ctx := context.Background()
	filter := bson.M{"_id": game.ID}
	update := bson.M{"$set": bson.M{
		"EndTime":   game.EndTime,
		"ID":        game.ID,
		"ObjectSeq": game.ObjectSeq,
		"Score":     game.Score,
		"StartTime": game.StartTime,
		"UserID":    game.UserID,
	}}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update game: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("game not found with id: %s", game.ID)
	}

	fmt.Printf("Updated %d games\n", result.ModifiedCount)

	return nil
}

// GetGamesByUser retrieves all games for a specific user from MongoDB
func (repo *MongoDbGameRepository) GetGamesByUser(userID string) ([]*domain.Game, error) {
	fmt.Println("")
	fmt.Println("GetGamesByUser", userID)
	ctx := context.Background()
	filter := bson.M{"UserID": userID}

	cursor, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error fetching games: %w", err)
	}
	// defer cursor.Close(ctx)

	var games []*domain.Game
	if err = cursor.All(ctx, &games); err != nil {
		fmt.Printf("Error decoding games for user %s: %v\n", userID, err)
		return nil, fmt.Errorf("error decoding games: %w", err)
	}

	fmt.Printf("Fetched %d games for user %s\n", len(games), userID)

	sort := func(a, b *domain.Game) int {
		return cmp.Compare(b.EndTime.UnixMilli(), a.EndTime.UnixMilli())
	}
	slices.SortFunc(games, sort)

	fmt.Println("len(games)", len(games))
	fmt.Println(games)
	fmt.Println("")

	return games, nil
}
