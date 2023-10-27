package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/briandowns/spinner"
)

var loading = spinner.New(spinner.CharSets[26], 200*time.Millisecond)

func GetJson[T any](url string, loadingText string) T {
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
		CommonExit(err)
	}
	return data
}
