package main

import (
	"github.com/al8n/media-convert/commands"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func main() {
	if err := commands.Execute(); err != nil {
		logrus.Error(err.Error())
	}
}
