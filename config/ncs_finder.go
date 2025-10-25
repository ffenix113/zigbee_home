package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/ffenix113/zigbee_home/types"
)

// In the future we want to move to nrfutil doing all this mess.

type NCSLocation struct {
	SDKVersion       types.Semver
	ToolchainVersion types.Semver
	ToolchainPath    string
	SDKPath          string
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

// FindNCSLocation will return paths for nRF SDK and Zephyr toolchains.
//
// If toolchain of required version was not found - it will try to
// use default version from toolchain file, and if it is not present -
// latest available in list of toolchains.
//
// As such this function can return different toolchain version that
// was requested, and caller can check it by comparing to version
// returned in NCSLocation.
func FindNCSLocation(sdkBasePath, sdkVersion string) (NCSLocation, error) {
	toolchainsJson := toolchainConfigPath(sdkBasePath)

	configFile, err := os.Open(toolchainsJson)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("open toolchains.json file at %q: %w", toolchainsJson, err)
	}
	defer configFile.Close()

	var toolchainFile []toolchainTopLevelItem

	err = json.NewDecoder(configFile).Decode(&toolchainFile)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("decode toolchain file: %w", err)
	}

	if len(toolchainFile) == 0 {
		return NCSLocation{}, fmt.Errorf("toolchain file does not contain definitions")
	}

	first := toolchainFile[0]

	return providePaths(sdkBasePath, sdkVersion, first)
}

func providePaths(ncsBase, sdkVersion string, toolchainItem toolchainTopLevelItem) (NCSLocation, error) {
	// Select SDK
	sdks, err := listSDKs(ncsBase)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("list sdks: %w", err)
	}

	if len(sdks) == 0 {
		return NCSLocation{}, fmt.Errorf("no sdks found in ncs bas of %q", ncsBase)
	}

	sort.Slice(sdks, func(i, j int) bool {
		// Prefer SDKs that are Zigbee add-ons.
		if sdks[i].IsZigbeeAddOn && !sdks[j].IsZigbeeAddOn {
			return true
		}

		return sdks[i].Version.Compare(sdks[j].Version) >= 0
	})

	var semverSDKVersion types.Semver

	if sdkVersion != "" {
		semverSDKVersion, err = types.ParseSemver(sdkVersion)
		if err != nil {
			return NCSLocation{}, fmt.Errorf("requested sdk version %q is invalid: %w", sdkVersion, err)
		}
	}

	var selectedSDK SDKInfo
	for _, sdk := range sdks {
		// Versions v3.0.0 and later do not contain Zigbee SDK,
		// so special SDK should be used.
		if sdk.Version[0] >= 3 && !sdk.IsZigbeeAddOn {
			continue
		}

		if sdk.Version.Compare(semverSDKVersion) >= 0 {
			selectedSDK = sdk
			break
		}
	}

	if selectedSDK.Path == "" {
		return NCSLocation{}, fmt.Errorf("sdk with version of at least %s was not found", semverSDKVersion)
	}

	// Select toolchain
	versionToIdentifier, err := mapVersions(toolchainItem)
	if err != nil {
		return NCSLocation{}, fmt.Errorf("map toolchain items: %w", err)
	}

	if len(versionToIdentifier) == 0 {
		return NCSLocation{}, fmt.Errorf("no toolchain versions found in toolchain configuration path %q", toolchainConfigPath(ncsBase))
	}

	var availableVersions []types.Semver
	for version := range versionToIdentifier {
		availableVersions = append(availableVersions, version)
	}

	// Sort versions in increasing order
	sort.Slice(availableVersions, func(i, j int) bool {
		return availableVersions[i].Compare(availableVersions[j]) == -1
	})

	log.Printf("requested sdk version: %q, available toolchain versions: %v", selectedSDK.Version, availableVersions)

	// Get directly requested or default version.
	selectedToolchainVersion := selectedSDK.Version
	bundleID := versionToIdentifier[selectedSDK.Version]

	// If directly requested version is not present - try to use latest from the same minor version.
	if bundleID == "" {
		selectedToolchainVersion, err = selectVersion(selectedSDK.Version, availableVersions)
		if err != nil {
			return NCSLocation{}, fmt.Errorf("select version: %w", err)
		}

		bundleID = versionToIdentifier[selectedToolchainVersion]
	}

	if bundleID == "" {
		return NCSLocation{}, errors.New("required toolchain version was not found and no other suitable version is present")
	}

	return NCSLocation{
		SDKVersion:       selectedSDK.Version,
		ToolchainVersion: selectedToolchainVersion,
		ToolchainPath:    filepath.Join(ncsBase, "toolchains", bundleID),
		SDKPath:          filepath.Join(selectedSDK.Path, "zephyr"),
	}, nil
}

type SDKInfo struct {
	Path          string
	Version       types.Semver
	IsZigbeeAddOn bool
}

func listSDKs(ncsBase string) ([]SDKInfo, error) {
	entries, err := os.ReadDir(ncsBase)
	if err != nil {
		return nil, fmt.Errorf("list ncs base path: %w", err)
	}

	var sdks []SDKInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		versionFilePath := path.Join(ncsBase, entry.Name(), "nrf", "VERSION")
		zigbeeAddOnPath := path.Join(ncsBase, entry.Name(), "ncs-zigbee")

		bts, err := os.ReadFile(versionFilePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			return nil, fmt.Errorf("read sdk version file %q: %w", versionFilePath, err)
		}

		sdkVersion, err := types.ParseSemver(strings.TrimSpace(string(bts)))
		if err != nil {
			return nil, fmt.Errorf("parse sdk version from %q: %w", versionFilePath, err)
		}

		_, err = os.Stat(zigbeeAddOnPath)
		isZigbeeAddOn := err == nil

		sdks = append(sdks, SDKInfo{
			Path:          path.Join(ncsBase, entry.Name()),
			Version:       sdkVersion,
			IsZigbeeAddOn: isZigbeeAddOn,
		})
	}

	return sdks, nil
}

func mapVersions(toolchainItem toolchainTopLevelItem) (map[types.Semver]string, error) {
	mapped := make(map[types.Semver]string, len(toolchainItem.Toolchains))

	for _, toolchain := range toolchainItem.Toolchains {
		for _, version := range toolchain.NCSVersions {
			semver, err := types.ParseSemver(version)
			if err != nil {
				return nil, fmt.Errorf("invalid ncs toolchain version %q: %w", version, err)
			}

			mapped[semver] = toolchain.Identifier.BundleID
		}
	}

	return mapped, nil
}

func selectVersion(requested types.Semver, available []types.Semver) (types.Semver, error) {
	foundIdx := -1
	for i, availableVer := range available {
		// Find first available version that it not less than what we look foor.
		if availableVer.Compare(requested) > 0 {
			foundIdx = i
			break
		}
	}

	if foundIdx == -1 {
		return types.Semver{}, nil
	}

	return available[foundIdx], nil
}

func toolchainConfigPath(ncsBase string) string {
	return filepath.Join(ncsBase, "toolchains", "toolchains.json")
}
