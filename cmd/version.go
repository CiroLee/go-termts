package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const VERSION_SHORT = "Print the version number of go-termts"
const VERSION = "0.0.8"

func init() {
	versionCmd.Flags().BoolP("version", "v", false, VERSION_SHORT)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   VERSION_SHORT,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
