package main

import (
	"fmt"
	"os"

	"github.com/rodrigo-brito/gocity/handle"
	"github.com/rodrigo-brito/gocity/lib"
	"github.com/rodrigo-brito/gocity/utils"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	defaultPort = 4000
)

func main() {
	analyzer := handle.AnalyzerHandle{
		Cache:   lib.NewCache(),
		Storage: lib.NewStorage(),
	}

	app := cli.NewApp()
	app.Version = "1.0.0"
	app.Description = "Code City metaphor for visualizing Go source code in 3D"
	app.Author = "Rodrigo Brito (https://github.com/rodrigo-brito)"

	app.Commands = []cli.Command{
		{
			Name:        "server",
			Description: "Start a local server to analyze projects",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "port",
					Value: defaultPort,
					Usage: "Local server port",
				},
			},
			Action: func(c *cli.Context) error {
				port := c.Int("port")
				return analyzer.Serve(port)
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

				analyzer.SetProjectURL(url)
				return analyzer.Serve(defaultPort)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
