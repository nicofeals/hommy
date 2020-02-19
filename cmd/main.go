package main

import (
	"os"
	"strings"
	"time"

	"github.com/nicofeals/hommy/common/config"
	"github.com/urfave/cli"
)

var version string

func main() {
	app := buildCLI()
	_ = app.Run(os.Args)
}

func getEnvironment(c *cli.Context) string {
	return strings.TrimSpace(c.GlobalString("env"))
}

func getLogLevel(c *cli.Context) string {
	return strings.TrimSpace(c.GlobalString("level"))
}

func getLightsOffDelay(c *cli.Context) time.Duration {
	return c.Duration("lights-off-delay")
}

func buildCLI() *cli.App {
	app := cli.NewApp()
	app.Name = "hommy"
	app.Usage = "Smart Home application"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "env, e",
			Value:  "development",
			Usage:  "runtime environment",
			EnvVar: config.EnvKeyEnvironment,
		},
		cli.StringFlag{
			Name:   "level",
			Value:  "info",
			Usage:  "logging level",
			EnvVar: config.EnvKeyLogLevel,
		},
		cli.StringFlag{
			Name:   "pg-addr",
			Value:  "localhost:5432",
			Usage:  "postgres db address",
			EnvVar: config.EnvKeyPostgresAddr,
		},
		cli.StringFlag{
			Name:   "pg-user",
			Value:  "postgres",
			Usage:  "postgres db user",
			EnvVar: config.EnvKeyPostgresUser,
		},
		cli.StringFlag{
			Name:   "pg-password",
			Value:  "postgres",
			Usage:  "postgres db password",
			EnvVar: config.EnvKeyPostgresPassword,
		},
		cli.StringFlag{
			Name:   "pg-database",
			Value:  "postgres",
			Usage:  "postgres database",
			EnvVar: config.EnvKeyPostgresDatabase,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start hommy server",
			Subcommands: cli.Commands{
				{
					Name:  "motionsensor",
					Usage: "start motion sensor server",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "port, p",
							Value:  "8080",
							Usage:  "port to start the motion sensor server on",
							EnvVar: config.EnvKeyServerPort,
						},
						cli.DurationFlag{
							Name:   "lights-off-delay",
							Value:  5 * time.Minute,
							Usage:  "time after which light turns off if no one's in the room",
							EnvVar: config.EnvKeyLightsOffDelay,
						},
					},
					Action: func(c *cli.Context) {
						startMotionSensorServerHandler(c)
					},
				},
			},
		},
	}

	return app
}
