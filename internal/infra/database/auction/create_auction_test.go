package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestNewAuctionRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should return repository with auctions collection when database is valid", func(mt *mtest.T) {
		db := mt.Client.Database("testdb")
		repo := NewAuctionRepository(db)

		require.NotNil(mt, repo, "repository should not be nil")
		require.NotNil(mt, repo.Collection, "collection should not be nil")

		assert.Equal(mt, "auctions", repo.Collection.Name(), "expected collection name to be 'auctions'")
	})

	mt.Run("should panic when database is nil", func(mt *mtest.T) {
		assert.Panics(mt, func() {
			NewAuctionRepository(nil)
		}, "expected panic when database is nil")
	})
}

func TestCreateAuction(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("should return internal error when InsertOne fails", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Message: "insert error"}))
		repo := &AuctionRepository{Collection: mt.Coll}

		auction := &auction_entity.Auction{
			Id:          "1",
			ProductName: "Test Product",
			Category:    "Electronics",
			Description: "A test product",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(context.Background(), auction)
		assert.NotNil(mt, err)
		assert.IsType(mt, &internal_error.InternalError{}, err)
		assert.Equal(mt, "Error trying to insert auction", err.Message)
	})

	mt.Run("should not return error when InsertOne succeeds and context is cancelled", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		repo := &AuctionRepository{Collection: mt.Coll}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		os.Setenv("AUCTION_INTERVAL", "20ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		auction := &auction_entity.Auction{
			Id:          "2",
			ProductName: "Laptop",
			Category:    "Tech",
			Description: "Gaming laptop",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(ctx, auction)
		assert.Nil(mt, err)

		time.Sleep(30 * time.Millisecond)
	})

	mt.Run("should call closeAuction successfully when InsertOne and UpdateOne succeed", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		repo := &AuctionRepository{Collection: mt.Coll}

		os.Setenv("AUCTION_INTERVAL", "10ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		ctx := context.Background()
		auction := &auction_entity.Auction{
			Id:          "3",
			ProductName: "TV",
			Category:    "Electronics",
			Description: "Smart TV",
			Condition:   auction_entity.Used,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(ctx, auction)
		assert.Nil(mt, err)

		time.Sleep(20 * time.Millisecond)
	})

	mt.Run("should handle UpdateOne error gracefully when closeAuction fails", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{Message: "update error"}))

		repo := &AuctionRepository{Collection: mt.Coll}

		os.Setenv("AUCTION_INTERVAL", "10ms")
		defer os.Unsetenv("AUCTION_INTERVAL")

		ctx := context.Background()
		auction := &auction_entity.Auction{
			Id:          "4",
			ProductName: "Phone",
			Category:    "Tech",
			Description: "Smartphone",
			Condition:   auction_entity.New,
			Status:      auction_entity.Active,
			Timestamp:   time.Now(),
		}

		err := repo.CreateAuction(ctx, auction)
		assert.Nil(mt, err)

		time.Sleep(20 * time.Millisecond)
	})
}

func TestGetAuctionInterval(t *testing.T) {

	t.Run("should return the duration from environment when AUCTION_INTERVAL is valid", func(t *testing.T) {
		os.Setenv("AUCTION_INTERVAL", "45s")
		defer os.Unsetenv("AUCTION_INTERVAL")

		interval := getAuctioInterval()
		assert.Equal(t, 45*time.Second, interval)
	})

	t.Run("should return default duration when AUCTION_INTERVAL is invalid", func(t *testing.T) {
		os.Setenv("AUCTION_INTERVAL", "invalid")
		defer os.Unsetenv("AUCTION_INTERVAL")

		interval := getAuctioInterval()
		assert.Equal(t, 2*time.Minute, interval)
	})

	t.Run("should return default duration when AUCTION_INTERVAL is not set", func(t *testing.T) {
		os.Unsetenv("AUCTION_INTERVAL")

		interval := getAuctioInterval()
		assert.Equal(t, 2*time.Minute, interval)
	})
}
