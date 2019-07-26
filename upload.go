package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	APIDEVKEY = "d36573a9390e1f80571857b20937dc12"
	USERKEY   = "a77afd7e1aab1e408072ee6dd9abc50e"
	PASTEURL  = "https://pastebin.com/api/api_post.php"
)

type pastePayload struct {
	APIDevKey       string `json:"api_dev_key"`
	APIUserKey      string `json:"api_user_key"`
	APIOption       string `json:"api_option"`
	APIPasteCode    string `json:"api_paste_code"`
	APIPasteName    string `json:"api_paste_name"`
	APIPastePrivate string `json:"api_paste_private"`
	APIPasteFormat  string `json:"api_paste_format"`
}

func createPaste(text string) string {
	v := url.Values{}
	v.Add("api_dev_key", APIDEVKEY)
	v.Add("api_user_key", USERKEY)
	v.Add("api_option", "paste")
	v.Add("api_paste_code", text)
	buf := strings.NewReader(v.Encode())
	// json.NewEncoder(buf).Encode(body)
	resp, err := http.Post(PASTEURL, "application/x-www-form-urlencoded", buf)
	if err != nil {
		log.Println("paste url: ", err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("read body: ", err)
	}
	return string(respBody)
}
