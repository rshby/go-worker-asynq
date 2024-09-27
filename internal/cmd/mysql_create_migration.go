package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	createMigrationCmd = &cobra.Command{
		Use:   "create-migration [filename]",
		Short: "Create new database migration file",
		Long:  "Create new database migration file with a specific name",
		Args:  cobra.ExactArgs(1),
		Run:   mysqlCreateMigration,
	}

	folderPath = "internal/database/migrations/"
)

func init() {
	RootCmd.AddCommand(createMigrationCmd)
}

func mysqlCreateMigration(cmd *cobra.Command, args []string) {
	// get filename from input argument
	filename := args[0]
	if err := checkMigrationFolderExists(); err != nil {
		logrus.Fatalf("error getting folder info : %s", err.Error())
		return
	}

	migrationFile := []byte(`-- +migrate Up notransaction` + "\n\n" + `-- +migrate Down`)
	migrationFileName := fmt.Sprintf("%s%s_%s.sql", folderPath, createUniqueTime(), strings.ToLower(filename))

	if err := ioutil.WriteFile(migrationFileName, migrationFile, 0666); err != nil {
		logrus.Fatalf("Error creating file: %s", err.Error())
	}

	logrus.Info(migrationFileName + " created")
}

// createUniqueTime is function to create unique time
func createUniqueTime() string {
	now := time.Now()
	splitDate := strings.Split(now.Format("01/02/2006"), "/")
	newDate := splitDate[2] + splitDate[0] + splitDate[1]

	hr, min, sc := now.Clock()
	hour := strconv.Itoa(hr)
	minute := strconv.Itoa(min)
	sec := strconv.Itoa(sc)

	if len(hour) == 1 {
		hour = "0" + hour
	}

	if len(minute) == 1 {
		minute = "0" + minute
	}

	if len(sec) == 1 {
		sec = "0" + sec
	}

	return newDate + hour + minute + sec
}

// checkMigrationFolderExists is function to check folder migrations exists or not
func checkMigrationFolderExists() error {
	_, err := os.Stat(folderPath)

	// if migrations folder not found
	if os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		logrus.Infof("[%s] folder not found, want to create (Y/N)?", folderPath)

		readString, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if !strings.Contains(readString, "Y") {
			return errors.New("cancelled creating migration")
		}

		// create migrations folder
		if err = createMigrationFolder(); err != nil {
			return err
		}
	}

	// folder migrations found
	return nil
}

// createMigrationFolder is function to create migrations folder
func createMigrationFolder() error {
	// cek folder database exists or not
	_, err := os.Stat("internal/database/")

	// if folder database not exists
	if os.IsNotExist(err) {
		// create folder database
		if err = os.Mkdir("database/", os.ModePerm); err != nil {
			logrus.Error(err)
			return err
		}
	}

	// create folder migrations
	if err = os.MkdirAll(folderPath, os.ModePerm); err != nil {
		logrus.Error(err)
		return err
	}

	// success create migrations folder
	return nil
}
