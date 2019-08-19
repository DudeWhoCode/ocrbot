package main

import (
	"github.com/otiai10/gosseract"
	"github.com/pkg/errors"
)

// read reads the image from given path and returns the transcribed text
func read(path string) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	err := client.SetImage(path)
	if err != nil {
		return "", errors.Wrap(err, "image read failed")
	}

	text, err := client.Text()
	if err != nil {
		return "", errors.Wrap(err, "text conversion failed")
	}

	return text, nil
}
