package main

import (
	"context"
	"os"

	"github.com/ffenix113/zigbee_home/cli/cmd/zigbee/firmware"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "zigbee",
		Usage: "Zigbee Home CLI application",
		Commands: []*cli.Command{
			firmware.RootCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "zigbee.yml",
			},
		},
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		panic(err)
	}
}
