package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/rodrigo-brito/gocity/handle"
	"github.com/rodrigo-brito/gocity/handle/middlewares"
	"github.com/rodrigo-brito/gocity/lib"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const EnvKeyGCS = "GOOGLE_APPLICATION_CREDENTIALS"

func main() {
	storage := lib.Storage(new(lib.NoStorage))
	router := chi.NewRouter()
	cache := lib.NewCache()

	// Use Google Cloud Storage for cache, if credentials available
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

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Description = "Code City metaphor for visualizing Go source code in 3D"
	app.Author = "Rodrigo Brito"

	app.Commands = []cli.Command{
		{
			Name:        "server",
			Description: "Start a local server to analyze projects",
			Action: func(c *cli.Context) error {
				router.Get("/api", analyzer.Handler)
				router.Get("/health", handle.HealthCheck)

				log.Println("Server started at http://localhost:4000")

				return http.ListenAndServe(":4000", router)
			},
		},
		{
			Name:        "open",
			Description: "Open a given project in local server",
			Action: func(c *cli.Context) error {
				fmt.Println(c.Args().First())
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
