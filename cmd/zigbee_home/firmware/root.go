package firmware

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

func RootCmd() *cli.Command {
	return &cli.Command{
		Name:  "firmware",
		Usage: "firmware operations like build & flash",
		Commands: []*cli.Command{
			buildCmd(),
			flashCmd(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "zigbee.yaml",
			},
			&cli.StringFlag{
				Name:  "workdir",
				Usage: "Change the working directory for the generation/build process. If this directory does not exist - it will be created.",
			},
		},
	}
}

func configFileFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Value: "zigbee.yaml",
		},
		&cli.StringSliceFlag{
			Name:    "override",
			Aliases: []string{"o"},
			Sources: cli.EnvVars("ZBHOME_OVERRIDE"),
			Usage:   "Takes one or more config options that will overwrite the one defined in configuration file. Also supports env variables in values.",
		},
	}
}

func parseConfig(configPath string, overwrites []string) (*config.Device, error) {
	cfg, err := parseConfigFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err = applyOverride(cfg, overwrites); err != nil {
		return nil, fmt.Errorf("apply overwrites: %w", err)
	}

	return cfg, nil
}

func parseConfigFile(configPath string) (*config.Device, error) {
	if configPath == "" {
		return nil, errors.New("config path cannot be empty (it is set by default)")
	}

	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("make config path %q absolute: %w", configPath, err)
	}

	conf, err := config.ParseFromFile(absConfigPath)
	if err != nil {
		return nil, fmt.Errorf("parse config file at %q: %w", absConfigPath, err)
	}

	return conf, nil
}

func applyOverride(cfg *config.Device, overwrites []string) error {
	var sb strings.Builder

	for _, overwrite := range overwrites {
		sb.Reset()

		separatorIdx := strings.IndexByte(overwrite, '=')
		if separatorIdx == -1 {
			return fmt.Errorf("overwrite %q does not have key and value separator ('=')", overwrite)
		}

		pathParts := strings.Split(overwrite[:separatorIdx], ".")
		for i, pathPart := range pathParts {
			if i != 0 {
				sb.WriteString(": { ")
			}

			sb.WriteString(pathPart)
		}

		sb.WriteString(": ")
		sb.WriteString(overwrite[separatorIdx+1:])
		sb.WriteString(strings.Repeat(" }", len(pathParts)-1))

		yamlString := os.ExpandEnv(strings.Replace(sb.String(), "=", ": ", 1))

		if err := yaml.Unmarshal([]byte(yamlString), cfg); err != nil {
			return fmt.Errorf("parse overwrite: %w", err)
		}
	}

	return nil
}
