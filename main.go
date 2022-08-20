package main

import (
	"fmt"
	"os"
	"time"

	"github.com/rodrigo-brito/gocity/pkg/lib"
	"github.com/rodrigo-brito/gocity/pkg/server"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	defaultPort = 4000
)

func main() {
	log.SetLevel(log.InfoLevel)

	app := cli.NewApp()

	app.Version = "1.0.6"
	app.Description = "Code City metaphor for visualizing Go source code in 3D"
	app.Copyright = "Rodrigo Brito (https://github.com/rodrigo-brito)"

	app.Commands = []*cli.Command{
		{
			Name:        "server",
			Description: "Start a local server to analyze projects",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "port",
					Aliases: []string{"b"},
					Value:   defaultPort,
					Usage:   "Local server port",
					EnvVars: []string{"PORT"},
				},
				&cli.DurationFlag{
					Name:    "cache",
					Aliases: []string{"c"},
					Value:   time.Hour,
					Usage:   "Cache's, TTL e.g.: --cache 4h",
					EnvVars: []string{"CACHE_TTL"},
				},
			},
			Action: func(c *cli.Context) error {
				analyzer := server.AnalyzerHandle{
					Cache:     lib.NewCache(),
					TmpFolder: os.TempDir(),
					CacheTTL:  c.Duration("cache"),
					Port:      c.Int("port"),
				}
				return analyzer.Serve()
			},
		},
		{
			Name:        "open",
			Description: "Open a given project in local server",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:    "port",
					Aliases: []string{"p"},
					Value:   defaultPort,
					Usage:   "Local server port",
					EnvVars: []string{"PORT"},
				},
				&cli.StringFlag{
					Name:    "branch",
					Aliases: []string{"b"},
					Value:   "master",
					Usage:   "Specify a custom branch",
				},
			},
			Action: func(c *cli.Context) error {
				var local bool
				projectAddress := c.Args().First()

				stat, err := os.Stat(projectAddress)
				if os.IsNotExist(err) {
					url, ok := lib.GetGithubBaseURL(c.Args().First())
					if !ok {
						return fmt.Errorf("project path not found")
					}
					projectAddress = url
				} else if !stat.IsDir() || projectAddress == "" {
					return fmt.Errorf("invalid project path")
				} else {
					local = true
				}

				analyzer := server.AnalyzerHandle{
					Cache:       lib.NewCache(),
					TmpFolder:   os.TempDir(),
					Port:        c.Int("port"),
					ProjectPath: &projectAddress,
					Local:       local,
				}

				if branch := c.String("branch"); branch != "" {
					analyzer.Branch = &branch
				}

				return analyzer.Serve()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
