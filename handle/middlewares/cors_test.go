package middlewares

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var GetCorsHeadersRequest = map[string]string {
	"Origin":                        "http://foobar.com",
	"Access-Control-Request-Method": "GET",
	"Access-Control-Request-Headers": "Accept, Authorization, Content-Type",
	"Access-Control-Allow-Methods": "*",
}

var GetCorsHeadersResponse = map[string]string {
	"Vary":                        "Origin",
	"Access-Control-Allow-Origin": "http://foobar.com",
	"Access-Control-Request-Headers": "Accept, Authorization, Content-Type",
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

func assertHeaders(t *testing.T, resHeaders http.Header, expHeaders map[string]string) {
	for _, name := range allHeaders {
		got := strings.Join(resHeaders[name], ", ")
		want := expHeaders[name]
		if got != want {
			t.Errorf("Response header %q = %q, want %q", name, got, want)
		}
	}
}

func TestGetCors(t *testing.T) {
	middleware := GetCors("*")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){})

	req , _ := http.NewRequest("GET","http://example.com/foo", nil)
	for name, value := range GetCorsHeadersRequest {
		req.Header.Add(name,value)
	}

	res := httptest.NewRecorder()
	middleware.Handler(handler).ServeHTTP(res, req)

	assertHeaders(t, res.Header(), GetCorsHeadersResponse)

}