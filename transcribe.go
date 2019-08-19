package main

import (
	"io"
	"net/http"
	"os"

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

// downaloadImage downloads the image from url to local file
func downloadImage(url, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "create image path failed")
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return errors.Wrap(err, "copy image to file failed")
	}

	return nil
}
