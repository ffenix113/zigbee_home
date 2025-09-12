package firmware

import (
	"github.com/urfave/cli/v2"
)

func RootCmd() *cli.Command {
	return &cli.Command{
		Name:  "firmware",
		Usage: "firmware operations like build & flash",
		Subcommands: []*cli.Command{
			buildCmd(),
			flashCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "workdir",
				Usage: "Change the working directory for the generation/build process. If this directory does not exist - it will be created.",
			},
		},
	}
}
