package release

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      18.12.2025
*/

type ReleaseType string

const (
	ReleaseMajor ReleaseType = "major"
	ReleaseMinor ReleaseType = "minor"
	ReleasePatch ReleaseType = "patch"
)

type Tool interface {
	Name() string
	Release(rt ReleaseType) error
}
