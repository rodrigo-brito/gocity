package handle

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/prometheus/common/log"
)

const frontendURL = "https://go-city.github.io"

func FrontEndProxy(w http.ResponseWriter, baseRequest *http.Request) {
	baseURL, err := url.Parse(frontendURL)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: baseURL.Scheme,
		Host:   baseURL.Host,
	})

	req := baseRequest
	req.URL.Scheme = baseURL.Scheme
	req.URL.Host = baseURL.Host
	req.Host = baseURL.Host

	proxy.ServeHTTP(w, req)
}
