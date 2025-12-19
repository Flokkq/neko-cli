package cmd

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      18.12.2025
*/

import (
	"strings"

	"github.com/nekoman-hq/neko-cli/internal/config"
	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/release"
	"github.com/spf13/cobra"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release [type]",
	Short: "Create a new release for your project",
	Long: `The release command helps you publish new versions of your project.
You can run it interactively or directly specify the type of release.

Examples:
  neko release          # starts an interactive survey to select release type
  neko release minor    # creates a minor release directly
  neko release major    # creates a major release directly
  neko release patch    # creates a patch release directly`,
	ValidArgs: []string{"major", "minor", "patch"},
	Args:      cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.LoadConfig()
		config.Validate(cfg)
		_ = config.GetPAT()

		tool, err := release.Get(string(cfg.ReleaseSystem))
		if err != nil {
			errors.Fatal(
				"Release System Not Found",
				err.Error(),
				errors.ErrInvalidReleaseSystem,
			)
			return
		}

		rt := release.ReleasePatch

		if len(args) > 0 {
			releaseType := strings.ToLower(args[0])

			switch releaseType {
			case "major":
				rt = release.ReleaseMajor
			case "minor":
				rt = release.ReleaseMinor
			case "patch":
				rt = release.ReleasePatch
			default:
				errors.Fatal(
					"Invalid Release Type",
					"Valid options are: major, minor, patch",
					errors.ErrInvalidReleaseType,
				)
				return
			}
		} else {
			// start Survey
			println("Start survey (no args found)")
			// TODO: Implement survey logic here
			// rt = startReleaseSurvey()
		}

		err = tool.Release(rt)
		if err != nil {
			errors.Fatal(
				"Release Failed",
				err.Error(),
				errors.ErrReleaseFailed,
			)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
}
