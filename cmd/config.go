package cmd

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/CiroLee/gear/gearmap"
	"github.com/CiroLee/go-termts/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

const GITHUB_URL = "https://api.github.com/repos/CiroLee/my-config/contents/"

var configMap = map[string]string{
	"prettier":   ".prettierrc",
	"commitlint": ".commitlintrc.js",
	"vscode":     ".vscode/settings.json",
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "download common used config files",
	Long:  "download common used config files, support prettier,commitlint,vscode(vscode-settings)",
	Run: func(cmd *cobra.Command, args []string) {
		download(args)
	},
}

func configToValue(config []string) []string {
	var result = make([]string, 0)
	for _, a := range config {
		for _, b := range gearmap.Keys(configMap) {
			if a == b {
				result = append(result, configMap[b])
				break
			}
		}
	}
	return result
}
func download(config []string) {
	r := configToValue(config)
	var name string
	for _, a := range r {
		url := GITHUB_URL + a
		data := utils.GetJson[map[string]any](url, "waiting...")
		content, err := base64.StdEncoding.DecodeString(data["content"].(string))
		if err != nil {
			utils.CommonExit(err)
		}
		wd, _ := os.Getwd()

		if strings.HasPrefix(a, ".vscode") {
			if !utils.Exists(wd + "/.vscode/") {
				os.Mkdir(wd+"/.vscode/", os.ModePerm)
				name = "/.vscode/settings.json"
			}
		} else {
			name = "/" + a
		}
		os.WriteFile(wd+name, content, os.ModePerm)
	}
}
