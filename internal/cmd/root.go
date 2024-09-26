package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-worker-asynq/internal/logger"
	"os"
)

var RootCmd = &cobra.Command{
	Use:   "go worker asynq",
	Short: "go worker asynq console",
	Long:  "this is go worker asynq console",
}

func init() {
	logger.SetupLogger()
}

// Execute is function to execute program
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
