package infrastructure

import (
	"context"
	"fmt"
	"os"

	"github.com/bernardbaker/qiba.core/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbReferralRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbReferralRepository() *MongoDbReferralRepository {
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
		fmt.Println("Referral repository - connection to MongoDB failed!")
	}
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Referral repository - Pinged your deployment. You successfully connected to MongoDB!")

	// if err := client.Database("admin").RunCommand(context.Background(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &MongoDbReferralRepository{
		client:     client,
		collection: client.Database("qiba-game").Collection("referrals"),
	}
}

// Save stores a new referral in MongoDB
func (repo *MongoDbReferralRepository) Save(obj *domain.Referral) error {
	ctx := context.Background()
	filter := bson.M{"_id": obj.ID}
	update := bson.M{"$set": obj}
	opts := options.Update().SetUpsert(true)

	_, err := repo.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to save referral: %w", err)
	}

	return nil
}

// Get retrieves a referral by its ID
func (repo *MongoDbReferralRepository) Get(objId string) *domain.Referral {
	ctx := context.Background()
	var referral domain.Referral
	filter := bson.M{"_id": objId}

	err := repo.collection.FindOne(ctx, filter).Decode(&referral)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf("referral not found with id: %s\n", objId)
			return nil
		}
		fmt.Printf("error fetching referral: %v\n", err)
		return nil
	}

	return &referral
}

// Update updates an existing referral in MongoDB
func (repo *MongoDbReferralRepository) Update(obj *domain.Referral) bool {
	ctx := context.Background()
	filter := bson.M{"_id": obj.ID}
	update := bson.M{"$set": obj}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		fmt.Printf("failed to update referral: %v\n", err)
		return false
	}

	if result.MatchedCount == 0 {
		fmt.Printf("referral not found with id: %s\n", obj.ID)
		return false
	}

	return true
}
