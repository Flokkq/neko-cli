package check

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      17.12.2025
*/

import (
	"regexp"

	"github.com/nekoman-hq/neko-cli/internal/config"
	"github.com/nekoman-hq/neko-cli/internal/errors"
)

var semverRegex = regexp.MustCompile(
	`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-[\da-zA-Z-]+(?:\.[\da-zA-Z-]+)*)?(?:\+[\da-zA-Z-]+(?:\.[\da-zA-Z-]+)*)?$`,
)

func ValidateConfig(cfg *config.NekoConfig) {
	if !cfg.ProjectType.IsValid() {
		errors.Error(
			"Invalid configuration",
			"ProjectType is invalid in .neko.json",
			errors.ErrConfigMarshal,
		)
		return
	}

	if !cfg.ReleaseSystem.IsValid() {
		errors.Error(
			"Invalid configuration",
			"ReleaseSystem is invalid in .neko.json",
			errors.ErrConfigMarshal,
		)
		return
	}

	if cfg.Version == "" {
		errors.Error(
			"Invalid configuration",
			"Version is missing in .neko.json",
			errors.ErrConfigMarshal,
		)
		return
	}

	if !semverRegex.MatchString(cfg.Version) {
		errors.Error(
			"Invalid configuration",
			"Version is not a valid semantic version (SemVer)",
			errors.ErrInvalidVersion,
		)
		return
	}

	println("\n✓ Configuration appears valid")
}
