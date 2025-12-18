package config

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      17.12.2025
*/

type ProjectType string
type ReleaseType string

const (
	ProjectTypeFrontend ProjectType = "frontend"
	ProjectTypeBackend  ProjectType = "backend"
	ProjectTypeOther    ProjectType = "other"
)

const (
	ReleaseTypeReleaseIt  ReleaseType = "release-it"
	ReleaseTypeJReleaser  ReleaseType = "jreleaser"
	ReleaseTypeGoReleaser ReleaseType = "goreleaser"
)

type NekoConfig struct {
	ProjectType   ProjectType `json:"projectType"`
	ReleaseSystem ReleaseType `json:"releaseSystem"`
	Version       string      `json:"version"`
}

func (p ProjectType) IsValid() bool {
	switch p {
	case ProjectTypeFrontend, ProjectTypeBackend, ProjectTypeOther:
		return true
	default:
		return false
	}
}

func (r ReleaseType) IsValid() bool {
	switch r {
	case ReleaseTypeReleaseIt, ReleaseTypeJReleaser, ReleaseTypeGoReleaser:
		return true
	default:
		return false
	}
}
