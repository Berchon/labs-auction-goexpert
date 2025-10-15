package mongodb

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Berchon/fullcycle-auction_go/configuration/logger"
	"github.com/Berchon/fullcycle-auction_go/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MONGODB_URL = "MONGODB_URL"
	MONGODB_DB  = "MONGODB_DB"
	APP_MODE    = "APP_MODE"
)

func NewMongoDBConnection(ctx context.Context) (*mongo.Database, error) {
	mongoURL := os.Getenv(MONGODB_URL)
	mongoDatabase := os.Getenv(MONGODB_DB)
	appMode := os.Getenv(APP_MODE)

	client, err := mongo.Connect(
		ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		logger.Error("Error trying to connect to mongodb database", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		logger.Error("Error trying to ping mongodb database", err)
		return nil, err
	}

	db := client.Database(mongoDatabase)

	if appMode == "dev" {
		err = ensureUsersCollection(ctx, db)
		if err != nil {
			return nil, err
		}
		err = ensureAuctionsCollection(ctx, db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func ensureUsersCollection(ctx context.Context, client *mongo.Database) error {
	collection := client.Collection("users")

	var count int64
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println("Error counting documents in users collection:", err)
		return err
	}

	if count == 0 {
		user := bson.M{
			"_id":  "e73fce6a-ccf5-4c12-87f7-30c5f9c9a6f7",
			"name": "user-test",
		}

		_, err := collection.InsertOne(ctx, user)
		if err != nil {
			log.Println("Error inserting user into users collection:", err)
			return err
		}
		log.Println("User inserted successfully into users collection")
	}

	return nil
}

func ensureAuctionsCollection(ctx context.Context, client *mongo.Database) error {
	collection := client.Collection("auctions")

	var count int64
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Println("Error counting documents in auctions collection:", err)
		return err
	}

	if count == 0 {
		id := "44c402b6-2960-4f9f-999f-5f217f40cee8"
		user := bson.M{
			"_id":          id,
			"product_name": "Mandolate",
			"category":     "Doce",
			"description":  "A melhor sobremesa do RU",
			"condition":    1,
			"status":       0,
			"timestamp":    time.Now().Unix(),
		}

		_, err := collection.InsertOne(ctx, user)
		if err != nil {
			log.Println("Error inserting auction into auctions collection:", err)
			return err
		}
		log.Println("Auction inserted successfully into auctions collection")

		go func() {
			select {
			case <-time.After(getAuctioInterval()):
				closeAuction(ctx, collection, id)
			case <-ctx.Done():
				logger.Error("Error to close auction, context cancelled", ctx.Err())
				return
			}
		}()
	}

	return nil
}

func closeAuction(ctx context.Context, collection *mongo.Collection, auctionId string) {
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
	filter := bson.M{"_id": auctionId}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("Error trying to update auction status to completed", err)
		return
	}
}

func getAuctioInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 2
	}

	return duration
}
