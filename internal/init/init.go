package init

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      23.12.2025
*/

import (
	"github.com/nekoman-hq/neko-cli/internal/config"
	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/git"
	"github.com/nekoman-hq/neko-cli/internal/release"
)

func Run(info *git.RepoInfo) {
	if !confirmOverwriteIfExists() {
		return
	}

	cfg := runWizard()

	if info != nil {
		cfg.ProjectOwner = info.Owner
		cfg.ProjectName = info.Repo
	}

	if err := config.SaveConfig(cfg); err != nil {
		errors.Fatal(
			"Configuration write failed",
			err.Error(),
			errors.ErrConfigWrite,
		)
		return
	}

	releaser, err := release.Get(string(cfg.ReleaseSystem))
	if err != nil {
		errors.Fatal(
			"Release System Not Found",
			err.Error(),
			errors.ErrInvalidReleaseSystem,
		)
	}

	err = releaser.Init(&cfg)
	if err != nil {
		errors.Fatal(
			"Release system initialization failed",
			"Failed to initialize the release system.",
			errors.ErrReleaseSystemInit,
		)
		return
	}

	printSetupInstructions(cfg)
}
