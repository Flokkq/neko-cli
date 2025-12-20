package release

import (
	"fmt"

	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/git"
)

func Preflight() {

	if err := git.IsClean(); err != nil {
		fmt.Println("⚠ Preflight checks failed!")
		errors.Error(
			"Uncommitted Changes Detected",
			err.Error(),
			errors.ErrDirtyWorkingTree,
		)
	}

	if err := git.OnMainBranch(); err != nil {
		fmt.Println("⚠ Preflight checks failed!")
		errors.Error(
			"Incorrect Git Branch",
			err.Error(),
			errors.ErrWrongBranch,
		)
	}
}
