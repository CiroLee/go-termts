package cmd

import (
	"fmt"
	"strings"

	"github.com/CiroLee/go-termts/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const BRANCH_COMMIT = "list all branches in current repo width an interactive way"

var remote bool

func init() {
	lsBranchCmd.Flags().BoolVarP(&remote, "remote", "r", false, BRANCH_COMMIT)
	rootCmd.AddCommand(lsBranchCmd)
}

var lsBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: BRANCH_COMMIT,
	Run: func(cmd *cobra.Command, args []string) {
		remote, _ := cmd.Flags().GetBool("remote")
		fmt.Printf("remote: %v\n", remote)
		url, err := utils.ExecuteCommand("git", "branch", "-r")
		if err != nil {
			utils.CommonExit(err)
		} else {
			branches := strings.Split(url, "\n")
			for i, v := range branches {
				branches[i] = strings.Trim(v, "* ")
			}
			prompt := promptui.Select{
				Label: "Select Branch",
				Items: branches,
				Size:  8,
			}
			_, b, err := prompt.Run()
			if err != nil {
				utils.CommonExit(err)
			} else {
				fmt.Printf("b: %v\n", b)
				// utils.ExecuteCommand("git", "checkout", b)
			}

		}
	},
}
