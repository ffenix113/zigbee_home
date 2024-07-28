package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/ffenix113/zigbee_home/types"
)

type NCSLocation struct {
	Version types.Semver
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

	if len(versionToIdentifier) == 0 {
		return NCSLocation{}, fmt.Errorf("no toolchain versions found in toolchain configuration path %q", toolchainConfigPath(ncsBase))
	}

	var availableVersions []types.Semver
	for version := range versionToIdentifier {
		parsed, err := types.ParseSemver(version)
		if err != nil {
			return NCSLocation{}, fmt.Errorf("parse semver %q: %w", version, err)
		}

		availableVersions = append(availableVersions, parsed)
	}

	// Sort versions in increasing order
	sort.Slice(availableVersions, func(i, j int) bool {
		return availableVersions[i].Compare(availableVersions[j]) == -1
	})

	log.Printf("requested toolchain version: %q, available versions: %v", version, availableVersions)

	// If version is not present - use default from the toolchain.
	// But if it is present - treat it as required.
	if version == "" {
		version = toolchainItem.DefaultToolchain["ncs_version"]
	}
	// Get directly requested or default version.
	bundleID := versionToIdentifier[version]

	// If directly requested version is not present - try to use latest from the same minor version.
	if bundleID == "" {
		requiredVersion, err := types.ParseSemver(version)
		if err != nil {
			return NCSLocation{}, fmt.Errorf("parse required version %q: %w", version, err)
		}

		foundVersion, err := selectVersion(requiredVersion, availableVersions)
		if err != nil {
			return NCSLocation{}, fmt.Errorf("select version: %w", err)
		}

		version = foundVersion.String()
		bundleID = versionToIdentifier[version]
	}

	if bundleID == "" {
		return NCSLocation{}, errors.New("required version was not found and no other suitable version is present")
	}

	semver, err := types.ParseSemver(version)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("parse found version %q: %w", version, err)
	}

	return constructPaths(ncsBase, semver, bundleID), nil
}

func mapVersions(toolchainItem toolchainTopLevelItem) map[string]string {
	mapped := make(map[string]string, len(toolchainItem.Toolchains))

	for _, toolchain := range toolchainItem.Toolchains {
		for _, version := range toolchain.NCSVersions {
			mapped[version] = toolchain.Identifier.BundleID
		}
	}

	return mapped
}

func selectVersion(requested types.Semver, available []types.Semver) (types.Semver, error) {
	foundIdx := -1
	for i, availableVer := range available {
		if availableVer.SameMajorMinor(requested) && availableVer.Compare(requested) > 0 {
			foundIdx = i
			requested = availableVer
		}
	}

	if foundIdx == -1 {
		return types.Semver{}, nil
	}

	return available[foundIdx], nil
}

func constructPaths(ncsBase string, version types.Semver, bundleID string) NCSLocation {
	return NCSLocation{
		Version: version,
		NCS:     filepath.Join(ncsBase, "toolchains", bundleID),
		Zephyr:  filepath.Join(ncsBase, version.String(), "zephyr"),
	}
}

func toolchainConfigPath(ncsBase string) string {
	return filepath.Join(ncsBase, "toolchains", "toolchains.json")
}
