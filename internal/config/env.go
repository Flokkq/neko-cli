package config

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      18.12.2025
*/

import (
	"fmt"
	"os"

	"github.com/nekoman-hq/neko-cli/internal/errors"
	"github.com/nekoman-hq/neko-cli/internal/log"
)

// GetPAT retrieves the GitHub Personal Access Token from the environment.
// It exits the program with a clear error message if the token is not set.
func GetPAT() string {
	log.V(log.Config, fmt.Sprintf("Looking up required env variable: %s",
		log.ColorText(log.ColorGreen, "GITHUB_TOKEN"),
	))
	token, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok || token == "" {
		errors.Fatal(
			"Environment Variable Missing",
			"A GitHub Personal Access Token (GITHUB_TOKEN) is required.\nSet it with: export GITHUB_TOKEN=your_token_here",
			errors.ErrMissingEnvVar,
		)
		// Fatal should exit, so the return is technically never reached
		return ""
	}
	return token
}
