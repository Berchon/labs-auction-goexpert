package dependencies

import (
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/database/auction"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/database/bid"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/database/user"
	"github.com/Berchon/fullcycle-auction_go/internal/usecase/auction_usecase"
	"github.com/Berchon/fullcycle-auction_go/internal/usecase/bid_usecase"
	"github.com/Berchon/fullcycle-auction_go/internal/usecase/user_usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func InitDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController) {

	auctionRepository := auction.NewAuctionRepository(database)
	bidRepository := bid.NewBidRepository(database, auctionRepository)
	userRepository := user.NewUserRepository(database)

	userController = user_controller.NewUserController(
		user_usecase.NewUserUseCase(userRepository))
	auctionController = auction_controller.NewAuctionController(
		auction_usecase.NewAuctionUseCase(auctionRepository, bidRepository))
	bidController = bid_controller.NewBidController(bid_usecase.NewBidUseCase(bidRepository))

	return
}
