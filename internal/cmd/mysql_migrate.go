package cmd

import (
	migrate "github.com/rubenv/sql-migrate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-worker-asynq/config"
	"go-worker-asynq/internal/database"
	"gorm.io/gorm"
	"strconv"
)

var (
	migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "migrate database",
		Long:  "command to migrate database",
		Run:   processMigration,
	}
)

func init() {
	migrateCmd.PersistentFlags().Int("step", 0, "maximum migration steps")
	migrateCmd.PersistentFlags().String("direction", "up", "migration direction")
	RootCmd.AddCommand(migrateCmd)
}

func processMigration(cmd *cobra.Command, args []string) {
	logrus.Info("running process migration")

	// get step and direction
	direction := cmd.Flag("direction").Value.String()
	stepStr := cmd.Flag("step").Value.String()
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		logrus.Fatalf("failed to migrate : %s", err.Error())
	}

	// open connection database mysql
	db, err := database.InitializeMySqlConnection()
	if err != nil {
		logrus.Fatalf("cant open connection to MYSQL [%s] : %s", config.MysqlDSN(), err.Error())
	}

	migration(db, direction, step)
}

func migration(db *gorm.DB, direction string, step int) {
	var (
		migrationDirection = migrate.Up
		migrations         = &migrate.FileMigrationSource{
			Dir: "internal/database/migrations/",
		}
	)

	migrate.SetTable("migrations")

	// get db mysql
	dbMysql, err := db.DB()
	if err != nil {
		logrus.Fatal(err)
	}
	defer dbMysql.Close()

	if direction == "down" {
		migrationDirection = migrate.Down
	}

	n, err := migrate.ExecMax(dbMysql, "mysql", migrations, migrationDirection, step)
	if err != nil {
		logrus.Fatalf("failed to migrate database : %s", err.Error())
	}

	logrus.Infof("Applied %d migrations!", n)
}
