package init

import (
	"encoding/json"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/nekoman-hq/neko-cli/internal/errors"
)

func Run() {
	// Check if config already exists
	if _, err := os.Stat(".neko.json"); err == nil {
		var overwrite bool
		confirm := &survey.Confirm{
			Message: ".neko.json already exists. Overwrite it?",
			Default: false,
		}

		if err := survey.AskOne(confirm, &overwrite); err != nil {
			errors.Warning(
				"Initialization cancelled",
				"Configuration wizard was aborted by the user.",
			)
			return
		}

		if !overwrite {
			errors.Warning(
				"Initialization aborted",
				"Existing .neko.json was not overwritten.",
			)
			return
		}
	}

	config := NekoConfig{}

	// Project type
	if err := survey.AskOne(&survey.Select{
		Message: "What kind of project is this?",
		Options: []string{"Frontend", "Backend", "Other"},
	}, &config.ProjectType); err != nil {
		errors.Error(
			"Project type selection failed",
			"Could not read project type input.",
			errors.ErrSurveyFailed,
		)
		return
	}

	// Release system
	if err := survey.AskOne(&survey.Select{
		Message: "Which release system should be used?",
		Options: getReleaseOptions(config.ProjectType),
	}, &config.ReleaseSystem); err != nil {
		errors.Error(
			"Release system selection failed",
			"Could not read release system input.",
			errors.ErrSurveyFailed,
		)
		return
	}

	// Initial version
	if err := survey.AskOne(&survey.Input{
		Message: "Initial version:",
		Default: "0.1.0",
		Help:    "Semantic Versioning (MAJOR.MINOR.PATCH)",
	}, &config.Version); err != nil {
		errors.Error(
			"Version input failed",
			"Could not read version input.",
			errors.ErrSurveyFailed,
		)
		return
	}

	if err := saveConfig(config); err != nil {
		errors.Fatal(
			"Configuration write failed",
			err.Error(),
			errors.ErrConfigWrite,
		)
		return
	}

	printSetupInstructions(config)
}

func getReleaseOptions(projectType string) []string {
	switch projectType {
	case "Frontend":
		return []string{"release-it"}
	case "Backend":
		return []string{"jreleaser"}
	case "Other":
		return []string{"goreleaser"}
	default:
		return []string{}
	}
}

func saveConfig(config NekoConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(".neko.json", data, 0644)
}

func printSetupInstructions(config NekoConfig) {
	println("\n✓ .neko.json created successfully\n")
	println("Next steps:")
	println("  • Use 'neko release' to create a release")
	println("  • Neko automatically manages the version in:")

	switch config.ReleaseSystem {
	case "release-it":
		println("    - package.json")
		println("    - .release-it.json")

	case "jreleaser":
		println("    - jreleaser.yml")
		println("    - pom.xml / build.gradle")

	case "goreleaser":
		println("    - .goreleaser.yml")
		println("    - Git tags")
	}

	println("\nTip: The version in .neko.json is the single source of truth.")
}
