package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(repoCmd)
}

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "open the repository of current project on your default browser",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := exec.Command("git config --get remote.origin.url").Output()
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		fmt.Printf("output: %v\n", output)

	},
}
