package router

import (
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/controller/auction_controller"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/controller/bid_controller"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/controller/user_controller"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	userController *user_controller.UserController,
	bidController *bid_controller.BidController,
	auctionController *auction_controller.AuctionController,
) {
	router.GET("/auction", auctionController.FindAuctions)
	router.GET("/auction/:auctionId", auctionController.FindAuctionById)
	router.POST("/auction", auctionController.CreateAuction)
	router.GET("/auction/winner/:auctionId", auctionController.FindWinningBidByAuctionId)

	router.POST("/bid", bidController.CreateBid)
	router.GET("/bid/:auctionId", bidController.FindBidByAuctionId)

	router.GET("/user/:userId", userController.FindUserById)
}
