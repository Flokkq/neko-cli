package release

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/log"
)

func ResolveReleaseType(args []string, t Tool, version *semver.Version) (Type, error) {
	if len(args) > 0 {
		rt, err := ParseReleaseType(args[0])
		if err != nil {
			errors.Fatal(
				"Not a valid increment",
				"The given type is not valid increment option.",
				errors.ErrInvalidReleaseType,
			)
		}

		newVer := NextVersion(version, rt)

		log.Print(log.Release,
			fmt.Sprintf(
				"Preview: current version %s → next version %s",
				log.ColorText(log.ColorCyan, version.String()),
				log.ColorText(log.ColorGreen, newVer.String()),
			),
		)

		return rt, nil
	}

	if !t.SupportsSurvey() {
		errors.Fatal(
			"Interactive mode not supported",
			fmt.Sprintf("%s requires an explicit release type", t.Name()),
			errors.ErrSurveyFailed,
		)
	}

	return t.Survey(version)
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
