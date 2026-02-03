// Package git includes operations using git or git-cli
package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/nekoman-hq/neko-cli/pkg/errors"
	"github.com/nekoman-hq/neko-cli/pkg/log"
)

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      20.12.2025
*/

func LatestTag() string {
	log.PluginV(log.Exec, fmt.Sprintf("%s (Extract last tag)", log.ColorText(log.ColorGreen, "git describe --tags --abbrev=0")))
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		errors.Warning(
			"Failed to get latest tag",
			"No tags found or could not execute git describe.\nUsing default version 0.1.0.",
		)
		return "0.1.0"
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		errors.Warning(
			"No tags found",
			"No tags exist in this repository.\nUsing default version 0.1.0.",
		)
		return "0.1.0"
	}

	log.PluginV(log.Guard, fmt.Sprintf("Latest tag: %s", outputStr))
	return outputStr
}

// GetTags returns a list of all git tags
func GetTags() []string {
	log.PluginV(log.Exec, "Fetching git tags: "+
		log.ColorText(log.ColorGreen, "git tag"))

	cmd := exec.Command("git", "tag")
	tagsOut, err := cmd.Output()
	if err != nil {
		errors.Warning(
			"Failed to fetch tags",
			fmt.Sprintf("Command failed: %s", err.Error()),
		)
		return []string{}
	}

	tagList := strings.Split(strings.TrimSpace(string(tagsOut)), "\n")
	if len(tagList) == 1 && tagList[0] == "" {
		return []string{}
	}

	return tagList
}

// CountCommitsBetween counts commits between two references
func CountCommitsBetween(from, to string) int {
	var cmd *exec.Cmd

	if from == "" {
		log.PluginV(log.Exec, fmt.Sprintf("Counting commits up to %s: %s",
			to, log.ColorText(log.ColorGreen, fmt.Sprintf("git rev-list --count %s", to))))
		cmd = exec.Command("git", "rev-list", "--count", to)
	} else {
		log.PluginV(log.Exec, fmt.Sprintf("Counting commits between %s and %s: %s",
			from, to, log.ColorText(log.ColorGreen, fmt.Sprintf("git rev-list --count %s..%s", from, to))))
		cmd = exec.Command("git", "rev-list", "--count", fmt.Sprintf("%s..%s", from, to))
	}

	out, err := cmd.Output()
	if err != nil {
		errors.Warning(
			"Failed to count commits",
			fmt.Sprintf("Command failed for range %s..%s: %s", from, to, err.Error()),
		)
		return 0
	}

	countStr := strings.TrimSpace(string(out))
	count, err := strconv.Atoi(countStr)
	if err != nil {
		errors.Warning(
			"Failed to parse commit count",
			fmt.Sprintf("Invalid count value: %s", countStr),
		)
		return 0
	}

	return count
}
