package main

import (
	"context"
	"log"

	"github.com/Berchon/fullcycle-auction_go/configuration/database/mongodb"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/router"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/dependencies"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
		return
	}

	databaseConnection, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	userController, bidController, auctionController := dependencies.InitDependencies(databaseConnection)

	r := gin.Default()
	router.RegisterRoutes(r, userController, bidController, auctionController)

	r.Run(":8080")
}
