package http_test

import (
	"net/http"
	"strings"
	"testing"

	http_test "github.com/Berchon/fullcycle-auction_go/tests/integration/http"
	"github.com/Berchon/fullcycle-auction_go/tests/integration/http/fixtures"
	"github.com/stretchr/testify/assert"
)

func TestCreateAuctionFailures(t *testing.T) {
	server := http_test.SetupServer(t)

	t.Run("should return 404 when body is missing", func(t *testing.T) {
		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", "")
		req.Header.Set("Content-Type", "application/json")

		resp := server.DoRequest(req)

		expected := `{"code":404,"err":"not_found","message":"Invalid type error","causes":null}`
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, expected, strings.TrimSpace(resp.Body.String()))

	})

	t.Run("should return 404 when JSON is malformed", func(t *testing.T) {
		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", `{"product_name": "Mola maluca", "category": "Brinquedo"`) // missing closing brace
		req.Header.Set("Content-Type", "application/json")

		resp := server.DoRequest(req)

		expected := `{"code":404,"err":"not_found","message":"Invalid type error","causes":null}`
		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, expected, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when required field is missing", func(t *testing.T) {
		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.MissingField)
		req.Header.Set("Content-Type", "application/json")

		resp := server.DoRequest(req)

		expected := `{
			"code":400,
			"err":"bad_request",
			"message":"Invalid field values",
			"causes":[{"field":"ProductName","message":"ProductName is a required field"}]
		}`
		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, expected, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 404 when field type is invalid", func(t *testing.T) {
		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidType)
		req.Header.Set("Content-Type", "application/json")

		resp := server.DoRequest(req)

		expected := `{
			"code":404,
			"err":"not_found",
			"message":"Invalid type error",
			"causes":null
		}`

		assert.Equal(t, http.StatusNotFound, resp.Code)
		assert.JSONEq(t, expected, strings.TrimSpace(resp.Body.String()))
	})

	t.Run("should return 400 when condition value is invalid", func(t *testing.T) {
		req := http_test.NewJSONRequest(t, http.MethodPost, "/auction", fixtures.InvalidCondition)
		req.Header.Set("Content-Type", "application/json")

		resp := server.DoRequest(req)

		expected := `{
			"code":400,
			"err":"bad_request",
			"message":"Invalid field values",
			"causes":[
					{"field":"Condition","message":"Condition must be one of [0 1 2]"}
			]
		}`

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		assert.JSONEq(t, expected, strings.TrimSpace(resp.Body.String()))
	})
}
