package release

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/nekoman-hq/neko-cli/pkg/log"
)

type Type string

const (
	Major Type = "major"
	Minor Type = "minor"
	Patch Type = "patch"
)

// ResolveReleaseType parses and validates the release type from the command argument
func ResolveReleaseType(version *semver.Version, releaseType Type) (Type, error) {
	newVer := NextVersion(version, releaseType)

	log.PluginPrint(log.Exec,
		"Applying %s (%s \uF178 %s)",
		log.ColorText(log.ColorPurple, string(releaseType)),
		version.String(),
		log.ColorText(log.ColorCyan, newVer.String()),
	)

	return releaseType, nil
}

func NextVersion(current *semver.Version, t Type) semver.Version {
	switch t {
	case Major:
		return current.IncMajor()
	case Minor:
		return current.IncMinor()
	case Patch:
		return current.IncPatch()
	default:
		return *current
	}
}

func ParseReleaseType(input string) (Type, error) {
	switch strings.ToLower(input) {
	case "major":
		return Major, nil
	case "minor":
		return Minor, nil
	case "patch":
		return Patch, nil
	default:
		// TODO - Handle Fatal Error
		return Patch, fmt.Errorf("valid options: major, minor, patch")
	}
}
