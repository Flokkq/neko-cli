package release

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

import (
	"github.com/nekoman-hq/neko-cli/pkg/errors"
	"github.com/nekoman-hq/neko-cli/pkg/log"
	"github.com/nekoman-hq/neko-cli/plugin/release/pkg/git"
)

func Preflight() {
	log.PluginV(log.Preflight, "Running pre-flight checks")
	if err := git.IsClean(); err != nil {
		errors.WriteError(
			"UNCOMMITTED_CHANGES",
			err.Error(),
		)
	}

	if err := git.EnsureNotDetached(); err != nil {
		errors.WriteError(
			"DETACHED_HEAD",
			err.Error(),
		)
	}

	if err := git.OnMainBranch(); err != nil {
		errors.WriteError(
			"INCORRECT_BRANCH",
			err.Error(),
		)
	}

	if err := git.HasUpstream(); err != nil {
		errors.WriteError(
			"NO_UPSTREAM_BRANCH",
			err.Error(),
		)
	}

	if err := git.IsUpToDate(); err != nil {
		errors.WriteError(
			"BRANCH_OUT_OF_DATE",
			err.Error(),
		)
	}

	log.PluginV(log.Preflight, "\uF00C Preflight checks succeeded!")
}
