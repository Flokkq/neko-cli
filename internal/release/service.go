package release

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

import (
	"fmt"

	"github.com/nekoman-hq/neko-cli/internal/config"
	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/git"
	"github.com/nekoman-hq/neko-cli/internal/log"
)

type Service struct {
	cfg *config.NekoConfig
}

func NewReleaseService(cfg *config.NekoConfig) *Service {
	return &Service{cfg: cfg}
}

func (rs *Service) Run(args []string) error {
	_, _ = git.Current()

	Preflight()
	VersionGuard(rs.cfg)

	releaser, err := Get(string(rs.cfg.ReleaseSystem))
	if err != nil {
		errors.Fatal(
			"Release System Not Found",
			err.Error(),
			errors.ErrInvalidReleaseSystem,
		)
	}

	rt, err := ResolveReleaseType(args, releaser)
	if err != nil {
		errors.Fatal(
			"Invalid Release Type",
			err.Error(),
			errors.ErrInvalidReleaseType,
		)
	}
	
	log.Print(log.VersionGuard, fmt.Sprintf("\uF00C All checks have succeeded. %s", log.ColorText(log.ColorGreen, "Starting release now!")))
	return releaser.Release(rt)
}
