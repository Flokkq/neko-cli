package config

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since     17.12.2025
*/

import (
	"encoding/json"
	"os"

	"github.com/nekoman-hq/neko-cli/internal/errors"
)

const configFileName = ".neko.json"

func LoadConfig() *NekoConfig {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		if os.IsNotExist(err) {
			errors.Fatal(
				"Configuration not found",
				"No .neko.json configuration found. Run 'neko init' first.",
				errors.ErrConfigNotExists,
			)
		} else {
			errors.Fatal(
				"Configuration read error",
				err.Error(),
				errors.ErrConfigRead,
			)
		}
	}

	var config NekoConfig
	if err := json.Unmarshal(data, &config); err != nil {
		errors.Fatal(
			"Configuration parse error",
			"Failed to parse .neko.json: "+err.Error(),
			errors.ErrConfigMarshal,
		)
	}

	return &config
}

func SaveConfig(config NekoConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		errors.Fatal(
			"Configuration serialization failed",
			"Could not marshal .neko.json: "+err.Error(),
			errors.ErrConfigMarshal,
		)
		return err
	}

	if err := os.WriteFile(configFileName, data, 0644); err != nil {
		errors.Fatal(
			"Configuration write failed",
			"Could not write .neko.json: "+err.Error(),
			errors.ErrConfigWrite,
		)
		return err
	}
	return nil
}
