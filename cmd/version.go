package cmd

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since     17.12.2025
*/

import (
	"fmt"
	"os"

	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show current version of this repo",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, ok := os.LookupEnv("GITHUB_TOKEN"); !ok {
			errors.Fatal(
				"Environment Variable Missing",
				"A Github Access Token (GITHUB_TOKEN) is required.\nSet it with: export GITHUB_TOKEN=your_token_here",
				errors.ErrMissingEnvVar,
			)
		}

		fmt.Println("version called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
