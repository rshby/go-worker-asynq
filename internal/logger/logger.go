package logger

import (
	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/sirupsen/logrus"
	"os"
)

// SetupLogger is function to setup logger
func SetupLogger() {
	formatter := runtime.Formatter{
		ChildFormatter: &logrus.TextFormatter{},
		Line:           true,
		File:           true,
	}

	logrus.SetFormatter(&formatter)
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}
