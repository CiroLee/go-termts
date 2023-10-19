package cmd

import (
	"regexp"
	"strings"

	"github.com/CiroLee/go-termts/utils"
	"github.com/spf13/cobra"
)

const REPO_SHORT = "open current git project in your default browser"

func processGitUrl(url string) string {
	url = strings.TrimSpace(url)
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = regexp.MustCompile(`^git@`).ReplaceAllString(url, "https://")
	}
	url = regexp.MustCompile(`\.git$`).ReplaceAllString(url, "")
	return url
}
func init() {
	repoCmd.Flags().BoolP("repo", "r", false, REPO_SHORT)
	rootCmd.AddCommand(repoCmd)
}

var repoCmd = &cobra.Command{
	Use:     "repo",
	Short:   REPO_SHORT,
	Aliases: []string{"r"},
	Run: func(cmd *cobra.Command, args []string) {
		url, err := utils.ExecuteCommand("git", "config", "--get", "remote.origin.url")
		if err != nil {
			utils.CommonExit(err)
		}
		url = processGitUrl(url)
		err = utils.Open(url)
		if err != nil {
			utils.CommonExit(err)
		}
	},
}
