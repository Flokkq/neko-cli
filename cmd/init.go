package cmd

/*
@Author     Benjamin Senekowitsch
@Contact    senekowitsch@nekoman.at
@Since     17.12.2025
*/

import (
	initcmd "github.com/nekoman-hq/neko-cli/internal/init"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize neko configuration",
	Long: `Interactive wizard to set up your project type and release system.
Neko manages version numbers uniformly across different release systems.`,
	Run: func(cmd *cobra.Command, args []string) {
		initcmd.Run()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
