package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/ffenix113/zigbee_home/cmd/zigbee_home/firmware"
	"github.com/urfave/cli/v3"
)

func main() {
	log.SetFlags(log.Lmsgprefix | log.LstdFlags | log.Lshortfile)

	appVer, goVer := getVersion()

	app := &cli.Command{
		Name:    "zigbee_home",
		Usage:   "Zigbee Home CLI application",
		Version: fmt.Sprintf("app: %s, built with: %s", appVer, goVer),
		Commands: []*cli.Command{
			firmware.RootCmd(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err.Error())
	}
}

func getVersion() (string, string) {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown", "unknown"
	}

	appVer := buildInfo.Main.Version
	if appVer == "" {
		appVer = "devel"
	}

	return appVer, buildInfo.GoVersion
}
