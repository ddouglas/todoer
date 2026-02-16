package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "todoer",
		Usage: "A todo list application ADHD-friendly reminder",
		Commands: []*cli.Command{
			serveCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logrus.WithError(err).Fatal("failed to run application")
	}

}
