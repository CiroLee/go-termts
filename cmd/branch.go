package cmd

import (
	"fmt"
	"strings"

	"github.com/CiroLee/gear/gearslice"
	"github.com/CiroLee/go-termts/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const BRANCH_COMMIT = "list branches and switch to target in current repo in an interactive way"

var remote bool

func init() {
	lsBranchCmd.Flags().BoolVarP(&remote, "remote", "r", false, "list remote branches")
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
			b = strings.Trim(b, "* ")
			b = strings.TrimPrefix(b, "origin/")
			if err != nil {
				utils.CommonExit(err)
			} else {
				switchBranch(b, remote)
			}
		}
	},
}

func switchBranch(tBranch string, remote bool) {
	if remote {
		// check if branch exists locally
		local, _ := utils.ExecuteCommand("git", "branch")
		localArr := strings.Split(local, "\n")
		localArr = utils.RemoveEmptyValues(localArr)
		if gearslice.Includes(localArr, tBranch) {
			utils.ExecuteCommand("git", "checkout", tBranch)
		} else {
			utils.ExecuteCommand("git", "checkout", "-b", tBranch)
		}
	} else {
		fmt.Println("Switching to local branch", tBranch)
		utils.ExecuteCommand("git", "checkout", tBranch)
	}

}
