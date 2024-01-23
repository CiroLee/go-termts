package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/CiroLee/go-termts/utils"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

type UserAndRepo struct {
	userName   string
	repository string
}

type RepoInfo struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

var branchName string
var dst string

func init() {
	dGitCmd.PersistentFlags().StringVar(&branchName, "branch", "", "branch name, default is repos's default branch")
	dGitCmd.PersistentFlags().StringVar(&dst, "dst", "", "destination directory, default is current directory")
	rootCmd.AddCommand(dGitCmd)
}

var loading = spinner.New(spinner.CharSets[26], 200*time.Millisecond)

var dGitCmd = &cobra.Command{
	Use:   "dgit",
	Short: "download github repository",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repoUrl := args[0]
		info, err := getUserAndRepo(repoUrl)
		if err != nil {
			log.Fatal(err)
		}

		if branchName == "" {
			branchName = getDefaultBranch(info.userName, info.repository)
		}
		if dst == "" {
			dst, _ = os.Getwd()
		}
		var downloadUrl = fmt.Sprintf("%s/archive/%s.tar.gz", repoUrl, branchName)
		err = downloadRepository(downloadUrl, dst, info.repository, branchName)
		if err != nil {
			log.Fatal("download failed: ", err)
		}
	},
}

func downloadRepository(gitUrl, dst, repository, branch string) error {
	loading.Prefix = "download"
	loading.Start()
	resp, err := http.Get(gitUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	unzipped, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(unzipped)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := header.Name
		branch = strings.ReplaceAll(branch, "/", "-")
		if strings.HasPrefix(name, fmt.Sprintf("%s-%s", repository, branch)) {
			name = strings.TrimPrefix(name, fmt.Sprintf("%s-%s/", repository, branch))
		}

		// process every file/folder in tar
		switch header.Typeflag {
		case tar.TypeDir:
			// create dir
			err := os.MkdirAll(filepath.Join(dst, name), 0755)
			if err != nil {
				return err
			}
		case tar.TypeReg:
			// write file
			outFile, err := os.Create(filepath.Join(dst, name))
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}

	loading.Stop()

	return nil

}

func getDefaultBranch(userName, repository string) string {
	gitUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s", userName, repository)
	resp := utils.GetJson[RepoInfo](gitUrl, "wait")

	return resp.DefaultBranch
}

func getUserAndRepo(repoUrl string) (UserAndRepo, error) {
	parseURL, err := url.ParseRequestURI(repoUrl)
	if err != nil {
		return UserAndRepo{}, err
	}
	if parseURL.Hostname() != "github.com" {
		log.Fatal("it is not an valid github repository url, please check")
	}
	rest := parseURL.Path[1:]
	rest = strings.Replace(rest, ".git", "", 1)
	userName := strings.Split(rest, "/")[0]
	repo := strings.Split(rest, "/")[1]
	return UserAndRepo{userName: userName, repository: repo}, nil
}
