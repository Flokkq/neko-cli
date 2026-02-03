package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose      bool
	outputFormat string
	pluginDir    string
	describe     bool
)

var rootCmd = &cobra.Command{
	Use:   "neko",
	Short: "Neko CLI - Plugin-based release and deployment tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Load plugins during initialization
	if err := InitializePlugins(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: Failed to initialize plugins: %v\n", err)
	}
}
