package goreleaser

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since      18.12.2025
*/

import (
	"fmt"

	"github.com/nekoman-hq/neko-cli/internal/release"
)

type Tool struct{}

func (t *Tool) Name() string {
	return "goreleaser"
}

func (t *Tool) Release(rt release.ReleaseType) error {
	fmt.Println("Goreleaser release:", rt)

	// 1. Version berechnen
	// 2. Git Tag erstellen
	// 3. Git Push
	// 4. goreleaser release ausführen

	return nil
}

func init() {
	release.Register(&Tool{})
}
