package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Berchon/fullcycle-auction_go/internal/infra/api/web/router"
	"github.com/Berchon/fullcycle-auction_go/internal/infra/dependencies"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type testServer struct {
	router *gin.Engine
	db     *mongo.Database
}

func SetupServer(t *testing.T, db *mongo.Database) *testServer {
	t.Helper()

	userController, bidController, auctionController := dependencies.InitDependencies(db)

	r := gin.Default()
	router.RegisterRoutes(r, userController, bidController, auctionController)

	s := &testServer{
		router: r,
		db:     db,
	}

	return s
}

func (s *testServer) DoRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.router.ServeHTTP(rr, req)
	return rr
}

func (s *testServer) DoRequestWithContext(ctx context.Context, req *http.Request) *httptest.ResponseRecorder {
	if ctx == nil {
		ctx = context.Background()
	}
	req = req.WithContext(ctx)
	return s.DoRequest(req)
}
