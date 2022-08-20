package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCors(t *testing.T) {

	var getCorsHeadersRequest = map[string]string{
		"Origin":                         "http://foobar.com",
		"Access-Control-Request-Method":  "GET",
		"Access-Control-Request-Headers": "Accept, Authorization, Content-Type",
		"Access-Control-Allow-Methods":   "*",
	}

	var getCorsHeadersResponse = map[string]string{
		"Vary":                             "Origin",
		"Access-Control-Allow-Origin":      "http://foobar.com",
		"Access-Control-Request-Headers":   "Accept, Authorization, Content-Type",
		"Access-Control-Allow-Credentials": "true",
	}

	var allHeaders = []string{
		"Vary",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Credentials",
		"Access-Control-Max-Age",
	}

	middleware := GetCors("http://foobar.com")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	req, _ := http.NewRequest("GET", "http://example.com/foo", nil)

	for name, value := range getCorsHeadersRequest {
		req.Header.Add(name, value)
	}

	res := httptest.NewRecorder()
	middleware.Handler(handler).ServeHTTP(res, req)

	for _, name := range allHeaders {
		require.Equal(t, getCorsHeadersResponse[name], strings.Join(res.Header()[name], ", "))
	}

}
