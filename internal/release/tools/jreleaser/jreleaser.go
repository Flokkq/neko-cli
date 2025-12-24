package jreleaser

import (
	"github.com/Masterminds/semver/v3"
	"github.com/nekoman-hq/neko-cli/internal/release"
)

type Jreleaser struct{}

func (j *Jreleaser) Name() string {
	return "jreleaser"
}

func (j *Jreleaser) Init(v *semver.Version) error {
	return nil
}

func (j *Jreleaser) Release(v *semver.Version) error {
	return nil
}

func (j *Jreleaser) Survey(v *semver.Version) (release.Type, error) {
	return release.Patch, nil
}

func (j *Jreleaser) SupportsSurvey() bool {
	return true
}
