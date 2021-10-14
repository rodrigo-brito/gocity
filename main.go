package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/rodrigo-brito/gocity/handle"
	"github.com/rodrigo-brito/gocity/lib"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	defaultPort = 4000
)

func main() {
	tmpFolder, err := ioutil.TempDir("", "")
	if err != nil {
		log.Fatal(err)
	}

	analyzer := handle.AnalyzerHandle{
		Cache:     lib.NewCache(),
		TmpFolder: tmpFolder,
	}

	log.SetLevel(log.InfoLevel)

	app := cli.NewApp()

	app.Version = "1.0.3"
	app.Description = "Code City metaphor for visualizing Go source code in 3D"
	app.Copyright = "Rodrigo Brito (https://github.com/rodrigo-brito)"

	app.Commands = []*cli.Command{
		{
			Name:        "server",
			Description: "Start a local server to analyze projects",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "port",
					Value:   defaultPort,
					Usage:   "Local server port",
					EnvVars: []string{"PORT"},
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
				url, ok := lib.GetGithubBaseURL(c.Args().First())
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
