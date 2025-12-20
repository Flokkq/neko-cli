package version

import (
	"fmt"
	"time"

	"github.com/nekoman-hq/neko-cli/internal/git"
	"github.com/nekoman-hq/neko-cli/internal/git/github"
)

var (
	// These variables are set via ldflags during build
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
	BuiltBy = "unknown"
)

func Latest(repoInfo *git.RepoInfo) {
	release := git.LatestRelease(repoInfo)
	displayCLIVersion()
	displayRelease(repoInfo, &release)
}

func displayCLIVersion() {
	fmt.Println()
	fmt.Printf("┌─ neko-cli\n")
	fmt.Printf("│\n")
	fmt.Printf("├─ Version:   %s\n", Version)
	fmt.Printf("├─ Commit:    %s\n", Commit)
	fmt.Printf("├─ Built:     %s\n", Date)
	fmt.Printf("└─ Built by:  %s\n", BuiltBy)
	fmt.Println()
}

func displayRelease(repoInfo *git.RepoInfo, release *github.Release) {
	// Parse and format the date
	publishedTime, err := time.Parse(time.RFC3339, release.PublishedAt)
	var formattedDate string
	if err == nil {
		formattedDate = publishedTime.Format("2006-01-02 15:04 MST")
	} else {
		formattedDate = release.PublishedAt
	}

	fmt.Println()
	fmt.Printf("┌─ Latest Release\n")
	fmt.Printf("│\n")
	fmt.Printf("├─ Repository: %s/%s\n", repoInfo.Owner, repoInfo.Repo)
	fmt.Printf("├─ Version:    %s", release.Name)
	if release.TagName != "" && release.TagName != release.Name {
		fmt.Printf(" (%s)", release.TagName)
	}
	fmt.Println()

	if release.PreRelease {
		fmt.Printf("├─ Type:       Pre-release\n")
	}

	fmt.Printf("├─ Published:  %s", formattedDate)
	if release.Author.Login != "" {
		fmt.Printf(" by %s", release.Author.Login)
	}
	fmt.Println()

	fmt.Printf("└─ URL:        %s\n", release.HTMLURL)
	fmt.Println()
}
