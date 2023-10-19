package utils

import (
	"log"
	"os"
	"os/exec"
)

func CommonExit(err error) {
	log.Println(err)
	os.Exit(1)
}

func ExecuteCommand(name string, subName string, args ...string) (string, error) {
	args = append([]string{subName}, args...)

	cmd := exec.Command(name, args...)
	bytes, err := cmd.CombinedOutput()

	return string(bytes), err
}
