package firmware

import (
	"os"
	"testing"

	"github.com/ffenix113/zigbee_home/config"
	"github.com/stretchr/testify/require"
)

func TestApplyOverrideValid(t *testing.T) {
	cfg := config.Device{
		General: config.General{
			Board: "replace me",
		},
	}

	overrides := []string{
		"general.board=one",
		"general.device_name=$USER",
	}

	require.NoError(t, applyOverride(&cfg, overrides))
	require.Equal(t, "one", cfg.General.Board)
	require.Equal(t, os.ExpandEnv("$USER"), cfg.General.DeviceName)
}

func TestApplyOverrideInvalid(t *testing.T) {
	cfg := config.Device{
		General: config.General{
			Board: "replace me",
		},
	}

	overrides := []string{
		"general.board22=one",
	}

	require.Error(t, applyOverride(&cfg, overrides))
}
