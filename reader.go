package main

import (
	"github.com/otiai10/gosseract"
)

func read(path string) string {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(path)
	text, _ := client.Text()
	return text
}
