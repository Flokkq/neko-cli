/*
Copyright © 2025 NAME HERE senekowitsch@nekoman.at
*/
package main

import (
	"github.com/nekoman-hq/neko-cli/cmd"
	_ "github.com/nekoman-hq/neko-cli/plugin/release/pkg/release/tool"
)

func main() {
	cmd.Execute()
}
