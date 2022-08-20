package middlewares

import (
	"net/http"
	"time"
)

const cookieKey = "gocity_api"

type Middleware func(http.Handler) http.Handler

func setAPICookie(w http.ResponseWriter, url string) {
	expire := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:    cookieKey,
		Value:   url,
		Path:    "/",
		Expires: expire,
		MaxAge:  90000,
	}
	http.SetCookie(w, &cookie)
}

func APIHeader(APIUrl string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			setAPICookie(w, APIUrl)
			next.ServeHTTP(w, r)
		})
	}
}
