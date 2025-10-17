package http_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var (
	MONGO_URI     string
	DATABASE_NAME string
)

func NewDB(t *testing.T) *DB {
	t.Helper()

	LoadEnv(t)

	if DATABASE_NAME == "auctions_db" {
		t.Fatal("🚨 Tests cannot run against the production database (auctions_db)!")
	}

	client := connectToMongoDB(t)
	db := client.Database(DATABASE_NAME)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	err := client.Ping(ctx, nil)
	require.NoError(t, err, "Failed to ping MongoDB after create connection")

	return &DB{
		Client:   client,
		Database: db,
	}
}

func LoadEnv(t *testing.T) {
	if err := godotenv.Load("../../cmd/auction/.env"); err != nil {
		t.Log("⚠️  Warning: .env file not found — using system environment variables")
	}

	originalMongoURI := os.Getenv("MONGODB_URL_TEST")
	originalDBName := os.Getenv("MONGODB_DB_TEST")

	require.NotEmpty(t, originalMongoURI, "❌ MONGODB_URL_TEST not set — please define it in your .env file")
	require.NotEmpty(t, originalDBName, "❌ MONGODB_DB_TEST not set — please define it in your .env file")

	timestamp := time.Now().UnixNano()
	newDBName := fmt.Sprintf("%s_%d", originalDBName, timestamp)

	// Replace the database name in the MongoDB URI with a temporary one that includes a timestamp.
	// By appending a unique timestamp to the database name, each test run uses an isolated database instance.
	// This ensures that integration tests do not interfere with each other, preventing race conditions or leftover data
	// from previous tests. Without this, using the same database name could lead to tests failing unpredictably
	// due to conflicts, data contamination, or concurrent modifications.
	MONGO_URI = strings.ReplaceAll(originalMongoURI, originalDBName, newDBName)
	DATABASE_NAME = newDBName
}

func connectToMongoDB(t *testing.T) *mongo.Client {
	t.Helper()

	clientOpts := options.Client().ApplyURI(MONGO_URI)
	client, err := mongo.Connect(context.Background(), clientOpts)
	require.NoError(t, err, "Failed to connect to MongoDB")
	return client
}

func (db *DB) Reconnect(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	err := db.Client.Ping(ctx, nil)
	if err != nil {
		client := connectToMongoDB(t)
		db.Client = client
		db.Database = client.Database(DATABASE_NAME)
	}

	err = db.Client.Ping(ctx, nil)
	require.NoError(t, err, "Failed to ping MongoDB after reconnection")
}

func (db *DB) Disconnect(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	if err := db.Client.Disconnect(ctx); err != nil {
		t.Logf("⚠️  Failed to disconnect MongoDB client: %v", err)
	}
}

func (db *DB) DropCollection(t *testing.T, collectionName string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	t.Cleanup(cancel)

	if err := db.Database.Collection(collectionName).Drop(ctx); err != nil {
		t.Logf("⚠️  Failed to drop collection %s: %v", collectionName, err)
	}
}

func (db *DB) DropAllCollections(t *testing.T) {
	t.Helper()

	collections := []string{"auctions", "bids", "users"}
	for _, collection := range collections {
		db.DropCollection(t, collection)
	}
}
