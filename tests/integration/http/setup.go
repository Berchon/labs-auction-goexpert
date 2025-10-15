package http_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/router"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/dependencies"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testServer struct {
	router *gin.Engine
	db     *mongo.Database
	client *mongo.Client
}

func init() {
	if err := godotenv.Load("../../cmd/auction/.env"); err != nil {
		log.Println("⚠️  Warning: .env file not found — using system environment variables")
	}

	if os.Getenv("MONGODB_URL_TEST") == "" {
		log.Fatal("❌ MONGODB_URL_TEST not set — please define it in your .env file")
	}
	if os.Getenv("MONGODB_DB_TEST") == "" {
		log.Fatal("❌ MONGODB_DB_TEST not set — please define it in your .env file")
	}
}

func SetupServer(t *testing.T) *testServer {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	mongoURI := os.Getenv("MONGODB_URL_TEST")
	dbName := os.Getenv("MONGODB_DB_TEST")

	require.NotEmpty(t, mongoURI, "MONGODB_URL_TEST must be set")
	require.NotEmpty(t, dbName, "MONGODB_DB_TEST must be set")

	if dbName == "auctions" {
		t.Fatal("🚨 Tests cannot run against the production database (auctions)!")
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOpts)
	require.NoError(t, err, "Failed to connect to MongoDB")

	db := client.Database(dbName)

	// Clean collections before each test
	for _, coll := range []string{"auctions", "bids", "users"} {
		if err := db.Collection(coll).Drop(ctx); err != nil {
			fmt.Printf("⚠️  Warning: failed to drop collection %s: %v\n", coll, err)
		}
	}

	userController, bidController, auctionController := dependencies.InitDependencies(db)

	r := gin.Default()
	router.RegisterRoutes(r, userController, bidController, auctionController)

	s := &testServer{
		router: r,
		db:     db,
		client: client,
	}

	// Ensure cleanup after test completion
	t.Cleanup(func() {
		s.Cleanup()
	})

	return s
}

func (s *testServer) DoRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	return rr
}

func (s *testServer) Cleanup() {
	if s.client != nil {
		if err := s.client.Disconnect(context.Background()); err != nil {
			log.Printf("⚠️  Warning: failed to disconnect MongoDB client: %v", err)
		}
	}
}
