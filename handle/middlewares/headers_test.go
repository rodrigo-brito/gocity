package middlewares

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const CookieHeader = "Set-Cookie"

func TestAPIHeader(t *testing.T) {
	now = func() time.Time {
		return time.Date(2020, time.October, 9, 1, 2, 3, 4, time.UTC)
	}

	url := "http://localhost:3000/test"
	expectedCookie := "gocity_api=http://localhost:3000/test; Path=/; Expires=Sat, 10 Oct 2020 01:02:03 GMT; Max-Age=90000"
	middleware := APIHeader(url)
	recorder := httptest.NewRecorder()

	router := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	router.ServeHTTP(recorder, nil)

	result := recorder.Result()
	assert.Equal(t, expectedCookie, result.Header.Get(CookieHeader))
}
