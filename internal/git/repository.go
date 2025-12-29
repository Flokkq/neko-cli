// Package git includes operations using git or git-cli
package git

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since     17.12.2025
*/

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/log"
)

type RepoInfo struct {
	Owner string
	Repo  string
}

type Contributor struct {
	Commits string
	Author  string
}

func Fetch() {
	log.V(log.VersionGuard, fmt.Sprintf("%s (Updating repository information)",
		log.ColorText(log.ColorGreen, "git fetch"),
	))

	exec.Command("git", "fetch")
}

// Current checks if a git repository exists and returns owner and repo name
func Current() (*RepoInfo, error) {
	log.V(log.Config, fmt.Sprintf("%s (Checking Repository Origin)",
		log.ColorText(log.ColorGreen, "git remote -v"),
	))

	cmd := exec.Command("git", "remote", "-v")
	output, err := cmd.CombinedOutput()
	if err != nil {
		errors.Fatal(
			"Not a Git Repository",
			"This directory is not a git repository.\nPlease run this command from within a git repository.",
			errors.ErrNoGitRepo,
		)
	}

	outputStr := string(output)
	if strings.TrimSpace(outputStr) == "" {
		errors.Fatal(
			"No Remote Found",
			"This git repository has no remote configured.\nAdd a remote with: git remote add origin <url>",
			errors.ErrNoRemote,
		)
	}
	return parseRemote(outputStr)
}

// parseRemote extracts owner and repo from git remote output
func parseRemote(remoteOutput string) (*RepoInfo, error) {
	// Regex patterns for both SSH and HTTPS URLs
	// SSH: git@git.com:owner/repo.git
	sshPattern := regexp.MustCompile(`git@github\.com:([^/]+)/([^/\s]+?)(?:\.git)?(?:\s|$)`)
	// HTTPS: https://github.com/owner/repo.git
	httpsPattern := regexp.MustCompile(`https://github\.com/([^/]+)/([^/\s]+?)(?:\.git)?(?:\s|$)`)

	// Try SSH pattern first
	if matches := sshPattern.FindStringSubmatch(remoteOutput); len(matches) >= 3 {
		repoPath := fmt.Sprintf("%s/%s", matches[1], matches[2])
		log.V(log.Config, fmt.Sprintf("Found repository: %s (SSH)",
			log.ColorText(log.ColorGreen, repoPath)))
		return &RepoInfo{
			Owner: matches[1],
			Repo:  matches[2],
		}, nil
	}

	// Try HTTPS pattern
	if matches := httpsPattern.FindStringSubmatch(remoteOutput); len(matches) >= 3 {
		repoPath := fmt.Sprintf("%s/%s", matches[1], matches[2])
		log.V(log.Config, fmt.Sprintf("Found repository: %s (HTTPS)",
			log.ColorText(log.ColorGreen, repoPath)))
		return &RepoInfo{
			Owner: matches[1],
			Repo:  matches[2],
		}, nil
	}

	errors.Fatal(
		"Invalid Remote URL",
		"Could not parse GitHub repository information from remote.\nOnly GitHub repositories are supported.",
		errors.ErrInvalidRemote,
	)

	return nil, nil // unreachable, but needed for compiler
}

func IsClean() error {
	log.V(log.Preflight, fmt.Sprintf("%s (Check branch state)",
		log.ColorText(log.ColorGreen, "git status --porcelain"),
	))
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to check git status: %w", err)
	}

	if strings.TrimSpace(string(output)) != "" {
		return fmt.Errorf("the working tree has uncommitted changes. Please commit or stash them")
	}

	log.V(log.Preflight, "Working tree is clean")
	return nil
}

func EnsureNotDetached() error {
	log.V(log.Preflight, fmt.Sprintf("%s (Ensure branch is not detached)",
		log.ColorText(log.ColorGreen, "git rev-parse --abbrev-ref HEAD"),
	))
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to determine HEAD state: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	if branch == "HEAD" {
		return fmt.Errorf("detached HEAD state detected. Please checkout a branch")
	}

	log.V(log.Preflight, fmt.Sprintf("HEAD attached to branch %s", log.ColorText(log.ColorGreen, branch)))
	return nil
}

func OnMainBranch() error {
	log.V(log.Preflight, fmt.Sprintf("%s (Check on main branch)",
		log.ColorText(log.ColorGreen, "git rev-parse --abbrev-ref HEAD"),
	))

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to determine current branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	if branch != "main" && branch != "master" {
		return fmt.Errorf("you are on branch '%s'. Releases are only allowed from 'main' or 'master'", branch)
	}

	log.V(log.Preflight, fmt.Sprintf("On %s branch", log.ColorText(log.ColorGreen, branch)))
	return nil
}

func HasUpstream() error {
	log.V(log.Preflight, fmt.Sprintf("%s (Check upstream configuration)",
		log.ColorText(log.ColorGreen, "git for-each-ref"),
	))

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to determine current branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))

	cmd = exec.Command(
		"git",
		"for-each-ref",
		"--format=%(upstream:short)",
		fmt.Sprintf("refs/heads/%s", branch),
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to determine upstream branch: %w", err)
	}

	upstream := strings.TrimSpace(string(output))
	if upstream == "" {
		return fmt.Errorf("branch '%s' has no upstream configured", branch)
	}

	log.V(log.Preflight, fmt.Sprintf("Upstream branch: %s", log.ColorText(log.ColorGreen, upstream)))
	return nil
}

func IsUpToDate() error {
	log.V(log.Preflight, fmt.Sprintf("%s (Check if branch is up to date)",
		log.ColorText(log.ColorGreen, "git status -sb"),
	))

	cmd := exec.Command("git", "status", "-sb")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("unable to check branch status: %w", err)
	}

	status := string(output)

	if strings.Contains(status, "behind") {
		return fmt.Errorf("branch is behind its upstream. Please pull the latest changes")
	}

	log.V(log.Preflight, "Branch is up to date with upstream")
	return nil
}

// CurrentBranch returns the name of the current branch
func CurrentBranch() string {
	log.V(log.History, "Fetching current branch: "+
		log.ColorText(log.ColorGreen, "git rev-parse --abbrev-ref HEAD"))

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchOut, err := cmd.Output()
	if err != nil {
		errors.Fatal(
			"Failed to get current branch",
			fmt.Sprintf("Command failed: %s", err.Error()),
			errors.ErrFileAccess,
		)
		return ""
	}

	branch := strings.TrimSpace(string(branchOut))
	return branch
}

// LastCommit returns the last commit information
func LastCommit() string {
	log.V(log.History, "Fetching last commit: "+
		log.ColorText(log.ColorGreen, "git log -1 --pretty=format:%%h '%%s' (%%cr)"))

	cmd := exec.Command("git", "log", "-1", "--pretty=format:%h '%s' (%cr)")
	lastCommitOut, err := cmd.Output()
	if err != nil {
		errors.Fatal(
			"Failed to get last commit",
			fmt.Sprintf("Command failed: %s", err.Error()),
			errors.ErrFileAccess,
		)
		return ""
	}

	lastCommit := strings.TrimSpace(string(lastCommitOut))
	return lastCommit
}

// TotalCommits returns the total number of commits as a string
func TotalCommits() string {
	log.V(log.History, "Counting total commits: "+
		log.ColorText(log.ColorGreen, "git rev-list --count HEAD"))

	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	totalCommitsOut, err := cmd.Output()
	if err != nil {
		errors.Warning(
			"Failed to count commits",
			fmt.Sprintf("Command failed: %s", err.Error()),
		)
		return "0"
	}

	return strings.TrimSpace(string(totalCommitsOut))
}

// FilesCount returns the number of tracked files
func FilesCount() int {
	log.V(log.History, "Counting tracked files: "+
		log.ColorText(log.ColorGreen, "git ls-files"))

	cmd := exec.Command("git", "ls-files")
	filesOut, err := cmd.Output()
	if err != nil {
		errors.Warning(
			"Failed to count files",
			fmt.Sprintf("Command failed: %s", err.Error()),
		)
		return 0
	}

	files := strings.Split(strings.TrimSpace(string(filesOut)), "\n")
	return len(files)
}

// RepoSize returns the repository size using du command
func RepoSize() string {
	log.V(log.History, "Calculating repository size: "+
		log.ColorText(log.ColorGreen, "du -sh ."))

	cmd := exec.Command("du", "-sh", ".")
	sizeOut, err := cmd.Output()
	if err != nil {
		log.V(log.History, "Could not determine repository size (du command not available)")
		return ""
	}

	fields := strings.Fields(string(sizeOut))
	if len(fields) == 0 {
		return ""
	}

	return fields[0]
}

// Contributors returns a list of contributors with their commit counts
func Contributors() []Contributor {
	log.V(log.History, "Fetching contributors: "+
		log.ColorText(log.ColorGreen, "git shortlog -sne HEAD"))

	cmd := exec.Command("git", "shortlog", "-sne", "HEAD")
	contrib, err := cmd.Output()
	if err != nil {
		errors.Fatal(
			"Failed to fetch contributors",
			fmt.Sprintf("Command failed: %s", err.Error()),
			errors.ErrFileAccess,
		)
		return []Contributor{}
	}

	contribLines := strings.Split(strings.TrimSpace(string(contrib)), "\n")
	log.V(log.History, fmt.Sprintf("Found %d contributors", len(contribLines)))

	contributors := make([]Contributor, 0, len(contribLines))
	for _, line := range contribLines {
		parts := strings.Fields(line)
		if len(parts) < 2 {
			log.V(log.History, fmt.Sprintf("Skipping invalid contributor line: %s", line))
			continue
		}

		contributors = append(contributors, Contributor{
			Commits: parts[0],
			Author:  strings.Join(parts[1:], " "),
		})
	}

	return contributors
}
