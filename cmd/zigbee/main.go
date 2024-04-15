package main

import (
	"context"
	"log"
	"os"

	"github.com/ffenix113/zigbee_home/cmd/zigbee/firmware"
	"github.com/urfave/cli/v2"
)

func main() {
	log.SetFlags(log.Lmsgprefix | log.LstdFlags | log.Lshortfile)

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
		log.Println(err.Error())
		os.Exit(1)
	}
}
