package utils

import (
	"log"
	"os"
	"os/exec"
	"runtime"
)

func CommonExit(err error) {
	log.Println(err)
	os.Exit(1)
}

func ExecuteCommand(name string, subName string, args ...string) (string, error) {
	args = append([]string{subName}, args...)
	args = RemoveEmptyValues(args)

	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()

	return string(bytes), err
}

func Open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func RemoveEmptyValues[T ~string](slice []T) []T {
	var result []T
	for _, v := range slice {
		if v != "" {
			result = append(result, v)
		}
	}
	return result
}
