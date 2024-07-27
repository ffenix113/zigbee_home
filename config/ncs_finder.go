package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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

	if len(versionToIdentifier) == 0 {
		return NCSLocation{}, fmt.Errorf("no toolchain versions found in toolchain configuration path %q", toolchainConfigPath(ncsBase))
	}

	var availableVersions []string
	for version := range versionToIdentifier {
		availableVersions = append(availableVersions, version)
	}

	sort.Strings(availableVersions)

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
		version = selectVersion(version, availableVersions)
		bundleID = versionToIdentifier[version]
	}

	if bundleID == "" {
		return NCSLocation{}, errors.New("required version was not found and no other suitable version is present")
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

type version [3]uint8

var versionRegx = regexp.MustCompile(`^v(\d+)\.(\d+)(?:\.(\d+))?$`)

func parseVersion(ver string) version {
	match := versionRegx.FindStringSubmatch(ver)
	if match == nil {
		log.Fatalf("incorrect version %q", ver)
	}

	var result version
	for i, part := range match[1:] {
		if i == 2 && part == "" {
			part = "0"
		}

		parsed, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			log.Fatalf("should not happen: bad part of the version: %q", part)
		}

		result[i] = uint8(parsed)
	}

	return result
}

func (v version) greaterOrEqual(other version) bool {
	return v[0] == other[0] &&
		v[1] == other[1] &&
		v[2] >= other[2]
}

func selectVersion(requested string, available []string) string {
	requestedVer := parseVersion(requested)

	foundIdx := -1
	for i, availableVer := range available {
		parsedAvailable := parseVersion(availableVer)

		if parsedAvailable.greaterOrEqual(requestedVer) {
			foundIdx = i
			requestedVer = parsedAvailable
		}
	}

	if foundIdx == -1 {
		return ""
	}

	return available[foundIdx]
}

func constructPaths(ncsBase, version, bundleID string) NCSLocation {
	return NCSLocation{
		Version: version,
		NCS:     filepath.Join(ncsBase, "toolchains", bundleID),
		Zephyr:  filepath.Join(ncsBase, version, "zephyr"),
	}
}

func toolchainConfigPath(ncsBase string) string {
	return filepath.Join(ncsBase, "toolchains", "toolchains.json")
}
