package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"slices"

	"golang.org/x/exp/maps"
)

type NCSLocation struct {
	Version string
	NCS     string
	Zephyr  string
}

type toolchainItem struct {
	Identifier struct {
		BundleID string `json:"bundle_id"`
	} `json:"identifier"`
	NCSVersions []string `json:"ncs_versions"`
}

type toolchainTopLevelItem struct {
	DefaultToolchain map[string]string `json:"default_toolchain"`
	Toolchains       []toolchainItem   `json:"toolchains"`
}

// FindNCSLocation will return paths for NCS and Zephyr toolchains.
//
// If toolchain of required version was not found - it will try to
// use default version from toolchain file, and if it is not present -
// latest available in list of toolchains.
//
// As such this function can return different toolchain version that
// was requested, and caller can check it by comparing to version
// returned in NCSLocation.
func FindNCSLocation(ncsBase, version string) (NCSLocation, error) {
	toolchainsJson := toolchainConfigPath(ncsBase)

	configFile, err := os.Open(toolchainsJson)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("open toolchains.json file at %q: %w", toolchainsJson, err)
	}

	var toolchainFile []toolchainTopLevelItem

	err = json.NewDecoder(configFile).Decode(&toolchainFile)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("decode toolchain file: %w", err)
	}

	if len(toolchainFile) == 0 {
		return NCSLocation{}, fmt.Errorf("toolchain file does not contain definitions")
	}

	first := toolchainFile[0]

	return providePaths(ncsBase, version, first)
}

func providePaths(ncsBase, version string, toolchainItem toolchainTopLevelItem) (NCSLocation, error) {
	versionToIdentifier := mapVersions(toolchainItem)

	availableVersions := maps.Keys(versionToIdentifier)
	if len(availableVersions) == 0 {
		return NCSLocation{}, fmt.Errorf("no toolchain versions found in toolchain configuration path %q", toolchainConfigPath(ncsBase))
	}

	log.Printf("available toolchain versions: %v", availableVersions)

	bundleID := versionToIdentifier[version]

	if bundleID == "" {
		version = toolchainItem.DefaultToolchain["ncs_version"]
		bundleID = versionToIdentifier[version]
	}

	if bundleID == "" {
		log.Println("toolchain config does not provide default version, and required version was not found, falling back to latest version in config")

		slices.Sort(availableVersions)
		version = availableVersions[len(availableVersions)-1]

		bundleID = versionToIdentifier[version]
	}

	return constructPaths(ncsBase, version, bundleID), nil
}

func mapVersions(toolchainItem toolchainTopLevelItem) map[string]string {
	mapped := make(map[string]string, len(toolchainItem.Toolchains))

	for _, toolchain := range toolchainItem.Toolchains {
		mapped[toolchain.NCSVersions[0]] = toolchain.Identifier.BundleID
	}

	return mapped
}

func constructPaths(ncsBase, version, bundleID string) NCSLocation {
	return NCSLocation{
		Version: version,
		NCS:     path.Join(ncsBase, "toolchains", bundleID),
		Zephyr:  path.Join(ncsBase, version, "zephyr"),
	}
}

func toolchainConfigPath(ncsBase string) string {
	return path.Join(ncsBase, "toolchains", "toolchains.json")
}
