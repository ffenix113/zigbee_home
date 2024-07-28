package examples

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/ffenix113/zigbee_home/cmd/zigbee/firmware"
	"github.com/ffenix113/zigbee_home/config"
	"github.com/stretchr/testify/require"
)

// This test will just ensure that the firmwares can be generated.
// It does not test that they are valid or can be built.
//
// Still, it is useful to see possible broken configurations,
// validate that templates and required functions.
func TestGenerateAllExamples(t *testing.T) {
	cwdAbsPath, err := filepath.Abs(".")
	require.NoError(t, err)

	examples, err := os.ReadDir(cwdAbsPath)
	require.NoError(t, err)

	for _, example := range examples {
		if !example.IsDir() {
			continue
		}

		t.Run(example.Name(), func(t *testing.T) {
			runExample(t, filepath.Join(cwdAbsPath, example.Name()))
		})
	}
}

// This test will generate the configuration from the repository root.
// Hopefuly resulting in less issues and up-to-date configuration.
func TestGenerateRootConfiguration(t *testing.T) {
	cwdAbsPath, err := filepath.Abs("./../")
	require.NoError(t, err)

	runExample(t, cwdAbsPath)
}

func runExample(t *testing.T, exampleDir string) {
	// FIXME: This is needed for config parser to
	// resolve includes in config files.
	os.Chdir(exampleDir)

	tmpDir, err := os.MkdirTemp("", "zigbee_home_*")
	require.NoError(t, err)

	t.Logf("created %s for example %s", tmpDir, exampleDir)

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	configPath := filepath.Join(exampleDir, "zigbee.yaml")
	// It is possible that examples would have some configuration
	// with non-default configuration file name.
	// If we will find this - skip such example.
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("skipping, as cannot access config file: %s", err.Error())
	}

	cfg, err := config.ParseFromFile(configPath)
	require.NoError(t, err)

	err = firmware.GenerateFirmwareFiles(context.Background(), tmpDir, false, cfg)
	require.NoError(t, err)
}
