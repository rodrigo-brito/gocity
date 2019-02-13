package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pkg/browser"

	"github.com/rodrigo-brito/gocity/utils"

	"github.com/rodrigo-brito/gocity/handle"
	"github.com/rodrigo-brito/gocity/lib"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	EnvKeyGCS   = "GOOGLE_APPLICATION_CREDENTIALS"
	defaultPort = 4000
)

func main() {
	storage := lib.Storage(new(lib.NoStorage))
	cache := lib.NewCache()

	if credentials := os.Getenv(EnvKeyGCS); len(credentials) > 0 {
		var err error
		storage, err = lib.NewGCS(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}

	analyzer := handle.AnalyzerHandle{
		Cache:   cache,
		Storage: storage,
	}

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Description = "Code City metaphor for visualizing Go source code in 3D"
	app.Author = "Rodrigo Brito (https://github.com/rodrigo-brito)"

	app.Commands = []cli.Command{
		{
			Name:        "server",
			Description: "Start a local server to analyze projects",
			Action: func(c *cli.Context) error {
				return analyzer.Serve(defaultPort)
			},
		},
		{
			Name:        "open",
			Description: "Open a given project in local server",
			Action: func(c *cli.Context) error {
				baseURL := c.Args().First()
				url, ok := utils.GetGithubBaseURL(c.Args().First())
				if !ok {
					return fmt.Errorf("invalid project URL: %s", baseURL)
				}
				go openVisualization(url, time.Second)
				return analyzer.Serve(defaultPort)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}

func openVisualization(url string, delay time.Duration) {
	time.Sleep(delay)
	cityURL := fmt.Sprintf("http://localhost:%d/#/%s", defaultPort, url)
	log.Info("opening ", cityURL)
	browser.OpenURL(cityURL)
}
