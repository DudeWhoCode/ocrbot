package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = *logrus.New()

func init() {
	log.Out = os.Stdout
	log.Level = logrus.InfoLevel
}
