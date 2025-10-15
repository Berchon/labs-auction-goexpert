package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func NewJSONRequest(t *testing.T, method, url string, body interface{}) *http.Request {
	t.Helper()

	var buf bytes.Buffer
	if body != nil {
		err := json.NewEncoder(&buf).Encode(body)
		require.NoError(t, err, "failed to encode JSON body")
	}

	req, err := http.NewRequest(method, url, &buf)
	require.NoError(t, err, "failed to create HTTP request")

	req.Header.Set("Content-Type", "application/json")
	return req
}
