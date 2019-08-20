package main

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var log = *logrus.New()

// newLogger creates a new logger instance that writes to both stdout and a log file
func newMultiLogger() (*os.File, error) {
	f, err := os.OpenFile("ocrbot.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "create file failed")
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Level = logrus.InfoLevel
	return f, nil
}
