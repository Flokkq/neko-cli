package releaseit

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      29.12.2025
*/

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Schema string        `json:"$schema"`
	Github GithubRelease `json:"github"`
	Git    GitConfig     `json:"git,omitempty"`
	Hooks  HooksConfig   `json:"hooks,omitempty"`
}

type GithubRelease struct {
	Release     bool   `json:"release"`
	ReleaseName string `json:"releaseName,omitempty"`
}

type GitConfig struct {
	Commit                 bool   `json:"commit"`
	Tag                    bool   `json:"tag"`
	Push                   bool   `json:"push"`
	RequireCleanWorkingDir bool   `json:"requireCleanWorkingDir"`
	Changelog              string `json:"changelog,omitempty"`
	CommitMessage          string `json:"commitMessage,omitempty"`
}

type HooksConfig struct {
	AfterBump string `json:"after:bump,omitempty"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(".release-it.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) (err error) {
	file, err := os.Create(".release-it.json")
	if err != nil {
		return fmt.Errorf("create .release-it.json: %w", err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("close file: %w", cerr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err = encoder.Encode(cfg); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	return nil
}

func InitDefaultConfig(projectName string) (*Config, error) {
	return &Config{
		Schema: "https://unpkg.com/release-it/schema/release-it.json",
		Github: GithubRelease{
			Release:     true,
			ReleaseName: fmt.Sprintf("%s@${version}", projectName),
		},
		Git: GitConfig{
			Commit:                 true,
			Tag:                    true,
			Push:                   true,
			RequireCleanWorkingDir: true,
			Changelog:              "npx auto-changelog --stdout --commit-limit false -u --template https://raw.githubusercontent.com/release-it/release-it/main/templates/changelog-compact.hbs",
			CommitMessage:          "chore(release): ${version}",
		},
		Hooks: HooksConfig{
			AfterBump: "npx auto-changelog -p",
		},
	}, nil
}
