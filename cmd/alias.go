package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(aliasCmd)
}

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "list the all alias from the .zshrc file",
	Run: func(cmd *cobra.Command, args []string) {
		p := getZshrcPath()
		file := readAliasFromZshrc(p)
		fmt.Println(strings.Join(file, ""))
	},
}

func getZshrcPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("error to get user home directory", err)
	}

	zshrcPath := fmt.Sprintf("%s/.zshrc", homeDir)
	return zshrcPath
}

func readAliasFromZshrc(filePath string) []string {
	var fileSlice []string
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("error to open the .zshrc file", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	blue := color.New(color.FgBlue).SprintFunc()
	white := color.New(color.FgWhite).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), "alias") {
			text := scanner.Text()
			short := strings.Replace(strings.Split(text, "\"")[0], "alias ", "", 1)
			cmd := strings.Split(text, short)[1]
			fileSlice = append(fileSlice, fmt.Sprintf("%s %s%s\n", blue("alias"), white(short), green(cmd)))
		}
	}

	return fileSlice
}
