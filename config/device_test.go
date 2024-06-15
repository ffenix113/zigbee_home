package config

import (
	"bytes"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigLoad(t *testing.T) {
	t.Run("valid small", func(t *testing.T) {
		data := []byte(`general: {board: some_board}`)

		cfg, err := ParseFromReader(&Device{}, bytes.NewReader(data))
		require.NoError(t, err)

		require.NotNil(t, cfg)
		require.Equal(t, "some_board", cfg.General.Board)
	})

	t.Run("invalid small", func(t *testing.T) {
		data := []byte(`some: {unknown: some_board}`)

		cfg, err := ParseFromReader(&Device{}, bytes.NewReader(data))
		require.Error(t, err)
		require.Nil(t, cfg)
	})
}

// TestLoadExamples will try to load example configurations and check if they are valid.
// It will not compile/generate them, only check if configuration is correct.
//
// In the future this should be replaced with actual building of the examples.
func TestLoadExamples(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	examplesPath := path.Join(cwd, "..", "examples")
	examples, err := os.ReadDir(examplesPath)
	if err != nil {
		t.Skipf("cannot read examples dir: %s", err.Error())
	}

	if len(examples) == 0 {
		t.Skip("no example directories")
	}

	for _, example := range examples {
		if !example.IsDir() {
			continue
		}

		t.Run(example.Name(), func(t *testing.T) {
			examplePath := path.Join(examplesPath, example.Name())
			// We need to change directory as some files can contain includes,
			// which will be relative to the configuration file.
			require.NoError(t, os.Chdir(examplePath))

			defaultConfig, err := os.Open(path.Join(examplePath, "zigbee.yaml"))
			if errors.Is(err, os.ErrNotExist) {
				t.Skip("default config not found")
			}

			require.NoError(t, err)

			defer defaultConfig.Close()

			_, err = ParseFromReader(&Device{}, defaultConfig)
			require.NoError(t, err)
		})
	}
}
