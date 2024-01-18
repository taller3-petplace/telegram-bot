package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"telegram-bot/src/app"
)

const logLevelEnv = "LOG_LEVEL"

func main() {
	logLevel := os.Getenv(logLevelEnv)
	if logLevel == "" {
		logLevel = "DEBUG"
	}

	err := initLogger(logLevel)
	if err != nil {
		fmt.Printf("error initializing logger: %v", err)
	}

	telegramer, err := app.NewTelegramer()
	if err != nil {
		logrus.Errorf("error creating telegramer: %v", err)
		return
	}

	logrus.Infoln("telegramer initialized correctly, lets get ready to rumble")
	err = telegramer.Start()
	logrus.Errorf("I'm gonna die %v", err)
}

// initLogger Receives the log level to be set in logrus as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func initLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	customFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   false,
	}
	logrus.SetFormatter(customFormatter)
	logrus.SetLevel(level)
	return nil
}
