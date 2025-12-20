package release

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

import "github.com/nekoman-hq/neko-cli/internal/git"

func VersionGuard() {
	// Git fetch
	git.Fetch()
	git.LatestTag()

	// Git latest tag
	// Compare version
}
