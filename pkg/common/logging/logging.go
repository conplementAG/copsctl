package logging

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	colorable "github.com/mattn/go-colorable"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	file          *os.File
	consoleLogger *logrus.Logger
	fileLogger    *logrus.Logger
)

func Initialize() {
	consoleLogger = logrus.New()
	fileLogger = logrus.New()

	fileLogger.SetLevel(logrus.DebugLevel)
	consoleLogger.SetLevel(logrus.InfoLevel)

	if viper.GetBool("verbose") {
		consoleLogger.SetLevel(logrus.DebugLevel)
	}

	consoleLogger.SetFormatter(&prefixed.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	fileLogger.SetFormatter(&prefixed.TextFormatter{
		DisableColors:   true,
		ForceFormatting: true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})

	file, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)

	if err != nil {
		panic(err)
	}

	consoleLogger.SetOutput(colorable.NewColorableStdout())
	fileLogger.SetOutput(file)
}

func Dispose() {
	file.Close()
}

func Info(text string) {
	consoleLogger.Info(text)
	fileLogger.Info(text)
}

func Infof(text string, v ...interface{}) {
	consoleLogger.Infof(text, v...)
	fileLogger.Infof(text, v...)
}

func Debug(text string) {
	consoleLogger.Debug(text)
	fileLogger.Debug(text)
}

func Debugf(text string, v ...interface{}) {
	consoleLogger.Debugf(text, v...)
	fileLogger.Debugf(text, v...)
}

func Fatalf(text string, v ...interface{}) {
	consoleLogger.Fatalf(text, v...)
	fileLogger.Fatalf(text, v...)
}

func Error(text string) {
	consoleLogger.Error(text)
	fileLogger.Error(text)
}

func Errorf(text string, v ...interface{}) {
	consoleLogger.Errorf(text, v...)
	fileLogger.Errorf(text, v...)
}
