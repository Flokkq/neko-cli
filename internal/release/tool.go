package release

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/nekoman-hq/neko-cli/internal/config"
	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/log"
)

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      18.12.2025
*/

type Tool interface {
	Name() string
	Init(cfg *config.NekoConfig) error
	Release(v *semver.Version) error
	Survey(v *semver.Version) (Type, error)
	SupportsSurvey() bool
}

type ToolBase struct{}

func (tb *ToolBase) RequireBinary(name string) {
	log.V(log.Init,
		fmt.Sprintf("Searching for %s executable: %s",
			name,
			log.ColorText(log.ColorGreen, fmt.Sprintf("which %s", name)),
		),
	)

	path, err := exec.LookPath(name)
	if err != nil {
		errors.Fatal(
			"Required dependency missing",
			fmt.Sprintf(
				"%s is not installed or not available in PATH",
				name,
			),
			errors.ErrDependencyMissing,
		)
	}

	log.Print(
		log.Init,
		"\uF00C Found %s at %s",
		log.ColorText(log.ColorCyan, name),
		log.ColorText(log.ColorGreen, path),
	)
}

// CreateReleaseCommit creates the chore commit for the release
func (tb *ToolBase) CreateReleaseCommit(v *semver.Version) error {
	commitMsg := fmt.Sprintf("chore(neko-release): %s", v)

	log.V(log.Release, fmt.Sprintf("Creating release commit: %s",
		log.ColorText(log.ColorGreen, fmt.Sprintf("git commit --allow-empty -m \"%s\"", commitMsg))))

	cmd := exec.Command("git", "commit", "--allow-empty", "-a", "-m", commitMsg)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errors.Fatal(
			"Failed to create release commit",
			fmt.Sprintf("git commit failed: %s", strings.TrimSpace(string(output))),
			errors.ErrReleaseCommit,
		)
	}

	log.Print(log.Release, "\uF00C Created release commit: %s",
		log.ColorText(log.ColorGreen, commitMsg))
	return nil
}

// CreateGitTag creates a git tag for the version
func (tb *ToolBase) CreateGitTag(v *semver.Version) error {
	tag := fmt.Sprintf("v%s", v)

	log.V(log.Release, fmt.Sprintf("Creating git tag: %s",
		log.ColorText(log.ColorGreen, fmt.Sprintf("git tag %s", tag))))

	cmd := exec.Command("git", "tag", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errors.Fatal(
			"Failed to create git tag",
			fmt.Sprintf("git tag %s failed: %s", tag, strings.TrimSpace(string(output))),
			errors.ErrReleaseTag,
		)
	}

	log.Print(log.Release, "\uF00C Created git tag: %s",
		log.ColorText(log.ColorGreen, tag))
	return nil
}

// PushCommits pushes the release commit to remote
func (tb *ToolBase) PushCommits() error {
	log.V(log.Release, fmt.Sprintf("Pushing release commit: %s",
		log.ColorText(log.ColorGreen, "git push origin HEAD")))

	cmd := exec.Command("git", "push", "origin", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		errors.Fatal(
			"Failed to push release commits",
			fmt.Sprintf("git push failed: %s", strings.TrimSpace(string(output))),
			errors.ErrReleasePush,
		)
	}

	log.Print(log.Release, "\uF00C Pushed release commit to %s",
		log.ColorText(log.ColorGreen, "origin"))
	return nil
}

// PushGitTag pushes the git tag to remote
func (tb *ToolBase) PushGitTag(v *semver.Version) error {
	tag := fmt.Sprintf("v%s", v)

	log.V(log.Release, fmt.Sprintf("Pushing git tag: %s",
		log.ColorText(log.ColorGreen, fmt.Sprintf("git push origin %s", tag))))

	cmd := exec.Command("git", "push", "origin", tag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		errors.Fatal(
			"Failed to push git tag",
			fmt.Sprintf("git push %s failed: %s", tag, strings.TrimSpace(string(output))),
			errors.ErrReleasePush,
		)
	}

	log.Print(log.Release, "\uF00C Pushed git tag: %s",
		log.ColorText(log.ColorGreen, tag))
	return nil
}
