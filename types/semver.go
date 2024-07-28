package types

import (
	"fmt"
	"regexp"
	"strconv"
)

type Semver [3]uint8

var versionRegx = regexp.MustCompile(`^v?(\d+)\.(\d+)(?:\.(\d+))?$`)

func NewSemver(major, minor, patch uint8) Semver {
	return Semver{major, minor, patch}
}

func ParseSemver(ver string) (Semver, error) {
	match := versionRegx.FindStringSubmatch(ver)
	if match == nil {
		return Semver{}, fmt.Errorf("incorrect version %q", ver)
	}

	parsePart := func(part string) (uint8, error) {
		if part == "" {
			return 0, nil
		}

		parsed, err := strconv.ParseUint(part, 10, 8)
		if err != nil {
			return 0, fmt.Errorf("should not happen: bad part of the version: %q", part)
		}

		return uint8(parsed), nil
	}

	var parsed [3]uint8
	for i, part := range match[1:] {
		uintPart, err := parsePart(part)
		if err != nil {
			return Semver{}, fmt.Errorf("parse part %q: %w", part, err)
		}
		parsed[i] = uintPart
	}

	return Semver{parsed[0], parsed[1], parsed[2]}, nil
}

func (s Semver) String() string {
	return fmt.Sprintf("v%d.%d.%d", s[0], s[1], s[2])
}

// SameMajorMinor checks if major and minor versions are equal.
func (s Semver) SameMajorMinor(another Semver) bool {
	return s[0] == another[0] && s[1] == another[1]
}

// Compare returns -1 if receiver is smaller than another,
// 1 if receiver is larger than another
// and 0 if they are equal.
func (s Semver) Compare(another Semver) int {
	if res := compare(s[0], another[0]); res != 0 {
		return res
	}

	if res := compare(s[1], another[1]); res != 0 {
		return res
	}

	return compare(s[2], another[2])
}

func compare(a, b uint8) int {
	if a > b {
		return 1
	}

	if a < b {
		return -1
	}

	return 0
}
