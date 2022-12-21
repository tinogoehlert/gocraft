package main

import (
	"os"

	"github.com/tinogoehlert/gocraft/internal/gocraft"
	cli "github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
		Name:        "GoCraft",
		Description: "a minecraft clone written in go ^^",
		Authors: []*cli.Author{{
			Name: "Tino Göhlert",
		}},
		DefaultCommand: "start",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:    "width",
				Value:   800,
				Aliases: []string{"sw"},
			},
			&cli.IntFlag{
				Name:    "height",
				Value:   450,
				Aliases: []string{"sh"},
			},
			&cli.StringFlag{
				Name:    "title",
				Value:   "GoCraft",
				Aliases: []string{"t"},
			},
			&cli.StringFlag{
				Name:    "loglevel",
				Value:   "info",
				Aliases: []string{"v"},
			},
			&cli.IntFlag{
				Name:  "fps",
				Value: 60,
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "start",
				Action: gocraft.RunEngine,
			},
		},
	}).Run(os.Args)
}
