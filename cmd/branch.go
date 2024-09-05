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
		var ext string
		if remote {
			ext = "-r"
		}
		branchesStr, err := utils.ExecuteCommand("git", "branch", ext)
		if err != nil {
			utils.CommonExit(err)
		} else {
			branches := strings.Split(branchesStr, "\n")
			branches = utils.RemoveEmptyValues(branches)
			prompt := promptui.Select{
				Label: "Select Branch",
				Items: branches,
				Size:  8,
			}
			_, b, err := prompt.Run()
			if err != nil {
				utils.CommonExit(err)
			} else {
				switchBranch(b, remote)
				// utils.ExecuteCommand("git", "checkout", b)
			}
		}
	},
}

func switchBranch(tBranch string, remote bool) {
	fmt.Printf("tBranch: %v\n", tBranch)
	fmt.Printf("remote: %v\n", remote)
	if remote {
		// check if branch exists locally
		local, _ := utils.ExecuteCommand("git", "branch")
		if strings.Contains(local, tBranch) {
			utils.ExecuteCommand("git", "checkout", tBranch)
		}
	} else {
		utils.ExecuteCommand("git", "checkout", "-b", tBranch)
	}

}
