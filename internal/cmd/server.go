package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-worker-asynq/internal/database"
)

var runServer = &cobra.Command{
	Use:   "server",
	Short: "run server",
	Long:  `This subcommand start the server`,
	Run:   server,
}

func init() {
	RootCmd.AddCommand(runServer)
}

func server(cmd *cobra.Command, args []string) {
	// connect database mysql
	db, err := database.InitializeMySqlConnection()
	if err != nil {
		logrus.Fatal(err)
	}

	dbMysql, err := db.DB()
	if err != nil {
		logrus.Fatal(err)
	}

	defer dbMysql.Close()
}
