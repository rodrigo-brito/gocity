package middlewares

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPIHeader(t *testing.T) {
	url := "http://localhost:3000/test"
	middleware := APIHeader(url)
	recorder := httptest.NewRecorder()

	router := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	router.ServeHTTP(recorder, nil)

	request := &http.Request{Header: http.Header{"Cookie": recorder.Header().Values("Set-Cookie")}}
	cookie, err := request.Cookie(cookieKey)
	require.NoError(t, err)
	require.Equal(t, cookie.Value, url)

	path, err := request.Cookie("Path")
	require.NoError(t, err)
	require.Equal(t, path.Value, "/")

	maxAge, err := request.Cookie("Max-Age")
	require.NoError(t, err)
	require.Equal(t, maxAge.Value, "90000")

	expiry, err := request.Cookie("Expires")
	require.NoError(t, err)

	expiryTime, err := time.Parse(time.RFC1123, expiry.Value)
	require.NoError(t, err)
	require.True(t, expiryTime.After(time.Now()))
}
