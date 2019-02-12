package main

import (
	"context"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"github.com/rodrigo-brito/gocity/handle"
	"github.com/rodrigo-brito/gocity/handle/middlewares"
	"github.com/rodrigo-brito/gocity/lib"
)

const EnvKeyGCS = "GOOGLE_APPLICATION_CREDENTIALS"

func main() {
	storage := lib.Storage(new(lib.NoStorage))
	router := chi.NewRouter()
	cache := lib.NewCache()

	// Use Google Cloud Storage for cache, if available
	if credentials := os.Getenv(EnvKeyGCS); len(credentials) > 0 {
		var err error
		storage, err = lib.NewGCS(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

	corsMiddleware := middlewares.GetCors("*")
	router.Use(corsMiddleware.Handler)

	analyzer := handle.AnalyzerHandle{
		Cache:   cache,
		Storage: storage,
	}

	router.Get("/api", analyzer.Handler)
	router.Get("/health", handle.HealthCheck)

	log.Println("Server started at http://localhost:4000")
	if err := http.ListenAndServe(":4000", router); err != nil {
		log.Error(err)
	}
}
