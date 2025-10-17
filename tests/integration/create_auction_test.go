package http_test

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Berchon/fullcycle-auction_go/internal/entity/auction_entity"
	http_test "github.com/Berchon/fullcycle-auction_go/tests/integration/http"
	"github.com/Berchon/fullcycle-auction_go/tests/integration/http/fixtures"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateAuction(t *testing.T) {
	// db := http_test.NewDB(t)
	// server := http_test.SetupServer(t, db.Database)

	// Validations in controller layer
	// -------------------------------------------------
	t.Run("should return 404 when body is missing", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", "")
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, http_test.InvalidTypeError, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 404 when JSON is malformed", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", `{"product_name": "Mola maluca", "category": "Brinquedo"`) // missing closing brace
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, http_test.MalformedJSONError, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when required field is missing", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.MissingField)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.ValidationError_ProductNameMissing, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 404 when field type is invalid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidType)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, http_test.InvalidTypeError, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when condition value is invalid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidCondition)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.ValidationError_InvalidCondition, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when product_name is invalid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidProductName)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.ValidationError_ProductNameMissing, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when description <= 10 and condition is invalid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidDescriptionAndCondition)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.ValidationError_DescriptionTooShort, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should not return error when description <= 10 but condition is valid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidShortDescription)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.ValidationError_DescriptionTooShort, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when multiple fields are invalid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.MultipleInvalidFields)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.ValidationError_DescriptionTooShort, strings.TrimSpace(resp.Body.String()))
	})

	// Validations in usecase and entity layers
	// -------------------------------------------------
	t.Run("should return 400 when category is invalid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidCategory)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, http_test.InvalidAuctionError, strings.TrimSpace(resp.Body.String()))
	})

	// Validations in repository
	// -------------------------------------------------
	t.Run("should return error when InsertOne fails and auction remains open in DB", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		// Simula falha no InsertOne
		err := db.Client.Disconnect(context.Background())
		assert.NoError(t, err, "failed to disconnect mongo client to simulate InsertOne failure")

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidAuction)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code, "should return 500 when insert fails")
		assert.JSONEq(t, http_test.InsertOneError, strings.TrimSpace(resp.Body.String()))

		db.Reconnect(t)

		coll := db.Database.Collection("auctions")
		var result auction_entity.Auction
		findErr := coll.FindOne(context.Background(), bson.M{"product_name": fixtures.ValidAuction["product_name"]}).Decode(&result)
		assert.Error(t, findErr, "auction should not exist in DB after failed InsertOne")
	})

	t.Run("should keep auction open when context is cancelled before closure", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		originalAppMode := os.Getenv("APP_MODE")
		os.Setenv("APP_MODE", "test")
		defer os.Setenv("APP_MODE", originalAppMode)

		originalInterval := os.Getenv("AUCTION_INTERVAL")
		os.Setenv("AUCTION_INTERVAL", "100ms")
		defer os.Setenv("AUCTION_INTERVAL", originalInterval)

		testCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidAuction)
		resp := server.DoRequestWithContext(testCtx, req)

		assert.Equal(t, http.StatusCreated, resp.Code, "should return 201 when auction is created successfully")
		assert.Empty(t, strings.TrimSpace(resp.Body.String()), "response body should be empty on successful creation")

		// Immediately cancels the context to simulate an early termination.
		// This mimics a scenario where the application shuts down or the HTTP request ends prematurely.
		cancel()

		// Waits for less than the auction closing interval (AUCTION_INTERVAL = 100ms).
		// The 50ms delay gives the goroutine enough time to detect the context cancellation
		// and handle it properly before the test checks the auction status.
		time.Sleep(50 * time.Millisecond)

		coll := db.Database.Collection("auctions")
		var result auction_entity.Auction
		findErr := coll.FindOne(context.Background(), bson.M{"product_name": fixtures.ValidAuction["product_name"]}).Decode(&result)
		assert.NoError(t, findErr, "auction should exist in DB after creation")
		assert.Equal(t, auction_entity.Active, result.Status, "auction should remain active before the interval expires")
	})

	t.Run("should log error when UpdateOne fails during automatic closure", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		originalInterval := os.Getenv("AUCTION_INTERVAL")
		os.Setenv("AUCTION_INTERVAL", "30ms")
		defer os.Setenv("AUCTION_INTERVAL", originalInterval)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidAuction)
		resp := server.DoRequest(req)
		assert.Equal(t, http.StatusCreated, resp.Code, "auction should be created successfully")
		assert.Empty(t, strings.TrimSpace(resp.Body.String()), "response body should be empty on successful creation")

		// Simulate a failure in the UpdateOne operation by disconnecting the MongoDB client
		// before the automatic auction closure can occur. This ensures that when the closure
		// routine tries to update the auction status in the database, it will fail.
		err := db.Client.Disconnect(context.Background())
		assert.NoError(t, err, "failed to disconnect MongoDB client to simulate UpdateOne failure")

		// Wait for the duration of the auction interval to allow the automatic closure
		// goroutine to attempt execution. This short pause ensures the test gives
		// enough time for the closure logic to run and encounter the simulated failure.
		time.Sleep(50 * time.Millisecond)

		// Reconnect the MongoDB client to validate the current status of the auction.
		// After reconnection, we can safely query the database to check if the auction
		// remains active despite the failed UpdateOne operation.
		db.Reconnect(t)
		coll := db.Database.Collection("auctions")

		var result auction_entity.Auction
		findErr := coll.FindOne(context.Background(), bson.M{"product_name": fixtures.ValidAuction["product_name"]}).Decode(&result)
		assert.NoError(t, findErr, "auction should still exist in DB")
		assert.Equal(t, auction_entity.Active, result.Status, "auction should remain active after failed UpdateOne")
	})

	t.Run("should keep auction open before the interval expires", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		originalAppMode := os.Getenv("APP_MODE")
		os.Setenv("APP_MODE", "test")
		defer os.Setenv("APP_MODE", originalAppMode)

		originalInterval := os.Getenv("AUCTION_INTERVAL")
		os.Setenv("AUCTION_INTERVAL", "50ms")
		defer os.Setenv("AUCTION_INTERVAL", originalInterval)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidAuction)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusCreated, resp.Code, "should return 201 when auction is created successfully")
		assert.Empty(t, strings.TrimSpace(resp.Body.String()), "response body should be empty on successful creation")

		coll := db.Database.Collection("auctions")
		var result auction_entity.Auction
		err := coll.FindOne(context.Background(), bson.M{"product_name": fixtures.ValidAuction["product_name"]}).Decode(&result)
		assert.NoError(t, err, "auction should exist in DB after creation")
		assert.Equal(t, auction_entity.Active, result.Status, "auction should remain active before the interval expires")
	})

	t.Run("should close multiple auctions concurrently", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		originalAppMode := os.Getenv("APP_MODE")
		os.Setenv("APP_MODE", "test")
		defer os.Setenv("APP_MODE", originalAppMode)

		originalInterval := os.Getenv("AUCTION_INTERVAL")
		os.Setenv("AUCTION_INTERVAL", "30ms")
		defer os.Setenv("AUCTION_INTERVAL", originalInterval)

		payloads := []map[string]interface{}{
			fixtures.ValidAuction,
			fixtures.ValidAuction2,
			fixtures.ValidAuction3,
		}

		var wg sync.WaitGroup
		for _, a := range payloads {
			wg.Add(1)
			go func(auction map[string]interface{}) {
				defer wg.Done()
				req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", auction)
				resp := server.DoRequest(req)
				assert.Equal(t, http.StatusCreated, resp.Code, "should return 201 when auction is created successfully")
				assert.Empty(t, strings.TrimSpace(resp.Body.String()), "response body should be empty on successful creation")
			}(a)
		}

		wg.Wait()

		time.Sleep(50 * time.Millisecond)

		coll := db.Database.Collection("auctions")
		cursor, err := coll.Find(context.Background(), bson.M{})
		assert.NoError(t, err, "should find all auctions")
		defer cursor.Close(context.Background())

		count := 0
		for cursor.Next(context.Background()) {
			var result auction_entity.Auction
			err := cursor.Decode(&result)
			assert.NoError(t, err)
			assert.Equal(t, auction_entity.Completed, result.Status, "auction should be closed after interval")
			count++
		}

		assert.Equal(t, len(payloads), count, "all created auctions should be found in DB")
	})

	t.Run("should automatically close the auction when interval expires", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		originalAppMode := os.Getenv("APP_MODE")
		os.Setenv("APP_MODE", "test")
		defer os.Setenv("APP_MODE", originalAppMode)

		originalInterval := os.Getenv("AUCTION_INTERVAL")
		os.Setenv("AUCTION_INTERVAL", "30ms")
		defer os.Setenv("AUCTION_INTERVAL", originalInterval)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidAuction)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusCreated, resp.Code, "should return 201 when auction is created successfully")
		assert.Empty(t, strings.TrimSpace(resp.Body.String()), "response body should be empty on successful creation")

		time.Sleep(50 * time.Millisecond)

		coll := db.Database.Collection("auctions")
		var result auction_entity.Auction
		err := coll.FindOne(context.Background(), bson.M{"product_name": fixtures.ValidAuction["product_name"]}).Decode(&result)
		assert.NoError(t, err, "auction should exist in DB")
		assert.Equal(t, auction_entity.Completed, result.Status, "auction should be automatically closed after interval")
	})

	t.Run("should create an auction successfully when payload is valid", func(t *testing.T) {
		db := http_test.NewDB(t)
		defer db.DropAllCollections(t)
		server := http_test.SetupServer(t, db.Database)

		originalAppMode := os.Getenv("APP_MODE")
		os.Setenv("APP_MODE", "test")
		defer os.Setenv("APP_MODE", originalAppMode)

		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.ValidAuction)
		resp := server.DoRequest(req)

		assert.Equal(t, http.StatusCreated, resp.Code, "should return 201 when auction is created successfully")
		assert.Empty(t, strings.TrimSpace(resp.Body.String()), "response body should be empty on successful creation")

		coll := db.Database.Collection("auctions")
		var result auction_entity.Auction
		err := coll.FindOne(context.Background(), bson.M{"product_name": fixtures.ValidAuction["product_name"]}).Decode(&result)
		assert.NoError(t, err, "auction should exist in DB after creation")
		assert.Equal(t, auction_entity.Active, result.Status, "auction status should be active after creation")
	})

}
