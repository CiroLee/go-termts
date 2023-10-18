package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CiroLee/gear/gearslice"
	"github.com/CiroLee/go-termts/utils"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const LICENSE_URL = "https://api.github.com/licenses"

var loading = spinner.New(spinner.CharSets[26], 200*time.Millisecond)

func init() {
	rootCmd.AddCommand(licenseCmd)
}

func getJson[T any](url string, loadingText string) T {
	loading.Prefix = loadingText
	loading.Start()
	r, err := http.Get(url)
	loading.Stop()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer r.Body.Close()
	var data T
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		utils.CommonExit(err)
	}
	return data
}

func getLicenseByKey(key string, license []map[string]string) string {
	r, ok := gearslice.Find(license, func(el map[string]string, _ int) bool {
		return el["key"] == key
	})
	if !ok {
		return ""
	}
	url := r["url"]
	data := getJson[map[string]any](url, "waiting...")
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
	Use:   "license",
	Short: "Output output LICENSE ",
	Run: func(cmd *cobra.Command, args []string) {
		// get license list
		license := getJson[[]map[string]string](LICENSE_URL, "loading")
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
					return errors.New("Invalid name")
				}
				return nil
			})
			if err != nil {
				utils.CommonExit(err)
			}
			year, err := inputPrompt("Year", func(input string) error {
				_, err := strconv.Atoi(input)
				if err != nil {
					return errors.New("Invalid year")
				}
				return nil
			})
			l = strings.Replace(l, "[fullname]", name, -1)
			l = strings.Replace(l, "[year]", year, -1)
		}
		os.WriteFile("LICENSE", []byte(l), 0644)
	},
}
