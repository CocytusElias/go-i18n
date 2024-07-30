package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:     "i18n bundle cli",
		Usage:    "goI18n command [command options]",
		Version:  "V1",
		HideHelp: true,
		Commands: []*cli.Command{
			genCmd,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
