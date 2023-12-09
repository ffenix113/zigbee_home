package firmware

import "github.com/urfave/cli/v2"

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
				Usage: "change the working directory for the build process (currently does not do anything)",
			},
		},
	}
}
