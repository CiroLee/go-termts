package utils

import (
	"log"
	"os"
)

func CommonExit(err error) {
	log.Println(err)
	os.Exit(1)
}
