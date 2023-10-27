package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CiroLee/gear/gearslice"
	"github.com/CiroLee/go-termts/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const LICENSE_URL = "https://api.github.com/licenses"
const LICENSE_SHORT = "Output output LICENSE"

func init() {
	licenseCmd.Flags().BoolP("license", "l", false, LICENSE_SHORT)
	rootCmd.AddCommand(licenseCmd)
}

func getLicenseByKey(key string, license []map[string]string) string {
	r, ok := gearslice.Find(license, func(el map[string]string, _ int) bool {
		return el["key"] == key
	})
	if !ok {
		return ""
	}
	url := r["url"]
	data := utils.GetJson[map[string]any](url, "waiting...")
	fmt.Printf("data: %T\n", data["body"])
	switch data["body"].(type) {
	case string:
		return data["body"].(string)
	}
	return ""
}

func needSignature(str string) bool {
	return strings.Contains(str, "[year]") && strings.Contains(str, "[fullname]")
}

func inputPrompt(label string, validate func(input string) error) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return result, nil
}

var licenseCmd = &cobra.Command{
	Use:     "license",
	Short:   LICENSE_SHORT,
	Aliases: []string{"l"},
	Run: func(cmd *cobra.Command, args []string) {
		// get license list
		license := utils.GetJson[[]map[string]string](LICENSE_URL, "loading")
		// select license key
		prompt := promptui.Select{
			Label: "Select a license",
			Items: gearslice.Map[map[string]string, string](license, func(el map[string]string, _ int) string {
				return el["key"]
			}),
			Size: 8,
		}
		_, key, err := prompt.Run()
		if err != nil {
			utils.CommonExit(err)
		}
		// get value by key
		l := getLicenseByKey(key, license)
		if needSignature(l) {
			name, err := inputPrompt("Name", func(input string) error {
				if len(input) < 1 {
					return errors.New("invalid name")
				}
				return nil
			})
			if err != nil {
				utils.CommonExit(err)
			}
			year, err := inputPrompt("Year", func(input string) error {
				_, err := strconv.Atoi(input)
				if err != nil {
					return errors.New("invalid year")
				}
				return nil
			})
			if err != nil {
				utils.CommonExit(err)
			}
			l = strings.Replace(l, "[fullname]", name, -1)
			l = strings.Replace(l, "[year]", year, -1)
		}
		os.WriteFile("LICENSE", []byte(l), 0644)
	},
}
