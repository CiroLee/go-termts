package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/CiroLee/gear/gearslice"
	"github.com/CiroLee/gear/gearstring"
	"github.com/CiroLee/go-termts/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

type CommitMsg struct {
	Type  string
	Value []string
}

var commitMsg = []CommitMsg{
	{
		Type: "feat",
		Value: []string{
			"(新增feature)",
			"(new feature)",
		},
	},
	{
		Type: "fix",
		Value: []string{
			"(修复bug)",
			"(fix bugs)",
		},
	},
	{
		Type: "style",
		Value: []string{
			"(仅修改了缩进、空格等样式，不改变代码逻辑)",
			"(Only the indentation and space styles have been changed, not the code logic)",
		},
	},
	{
		Type: "refactor",
		Value: []string{
			"(代码重构, 没有加新功能或者修复 bug)",
			"(Code refactoring, no new features or bug fixes)",
		},
	},
	{
		Type: "chore",
		Value: []string{
			"(改变构建流程、或者增加依赖库、工具等)",
			"(Change the build process, or add dependencies, tools, etc)",
		},
	},
	{
		Type: "docs",
		Value: []string{
			"(仅修改文档, 如README, CHANGELOG, CONTRIBUTE等)",
			"(Modification of documents only, e.g. README, CHANGELOG, CONTRIBUTE, etc)",
		},
	},
	{
		Type: "perf",
		Value: []string{
			"(优化相关, 比如提升性能、体验)",
			"(Optimization related, e.g. to improve performance, experience)",
		},
	},
	{
		Type: "revert",
		Value: []string{
			"(回滚到上一个版本)",
			"(Revert to the previous version)",
		},
	},
	{
		Type: "test",
		Value: []string{
			"(测试用例，包括单元测试，集成测试等)",
			"(Test cases, including unit tests, integration tests, etc.)",
		},
	},
}

func commitPrompts(lang string) []string {
	list := gearslice.Map[CommitMsg, string](commitMsg, func(el CommitMsg, _ int) string {
		var str string
		if lang == "zh" {
			str += el.Type + ": " + el.Value[0]
		} else {
			str += el.Type + ": " + el.Value[1]
		}
		return str
	})

	return list
}

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "shortcut for git commit, support zh(for Chinese) and en(for English) flags",
	Run: func(cmd *cobra.Command, args []string) {
		var lang string
		// default zh
		if len(args) == 0 || (len(args) > 0 && !gearslice.Includes([]string{"zh", "en"}, args[0])) {
			lang = "zh"
		}
		prompt := promptui.Select{
			Label: "Select a commit type",
			Items: commitPrompts(lang),
			Size:  8,
		}
		_, result, err := prompt.Run()
		if err != nil {
			utils.CommonExit(err)
		}
		gitType := strings.Split(result, ":")[0]
		msgPrompt := promptui.Prompt{
			Label: "please input commit message",
			Validate: func(input string) error {
				if len(input) < 1 {
					return errors.New("invalid message")
				}
				return nil
			},
		}
		msg, err := msgPrompt.Run()
		if err != nil {
			utils.CommonExit(err)
		}

		str, err := utils.ExecuteCommand("git", "commit", "-m", gearstring.Contact(`"`, gitType, ": ", msg, `"`))
		fmt.Println(str)
		if err != nil {
			utils.CommonExit(err)
		}
	},
}
