package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ffenix113/zigbee_home/cmd/zigbee/firmware"
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
				Value: "zigbee.yaml",
			},
			&cli.StringFlag{
				Name:    "log-level",
				EnvVars: []string{"ZBHOME_LOG_LEVEL"},
				Action: func(_ *cli.Context, s string) error {
					var lvl slog.Level
					if err := lvl.UnmarshalText([]byte(s)); err != nil {
						return fmt.Errorf("parse log level: %w", err)
					}

					slog.SetLogLoggerLevel(lvl)

					return nil
				},
			},
		},
	}

	if err := app.RunContext(context.Background(), os.Args); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
