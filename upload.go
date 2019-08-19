package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
)

var (
	apiDevKey   = os.Getenv("APIDEVKEY")
	userKey     = os.Getenv("USERKEY")
	pasteBinURL = "https://pastebin.com/api/api_post.php"
)

// createPaste uploads the text to pastebin account and returns URL of the paste
func createPaste(text string) (string, error) {
	v := url.Values{}
	v.Add("api_dev_key", apiDevKey)
	v.Add("api_user_key", userKey)
	v.Add("api_option", "paste")
	v.Add("api_paste_code", text)

	buf := strings.NewReader(v.Encode())
	resp, err := http.Post(pasteBinURL, "application/x-www-form-urlencoded", buf)
	if err != nil {
		return "", errors.Wrap(err, "pastebin upload failed")
	}

	pasteURL, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read failed")
	}
	return string(pasteURL), nil
}
