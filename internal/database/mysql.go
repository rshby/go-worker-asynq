package database

import (
	"github.com/sirupsen/logrus"
	"go-worker-asynq/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeMySqlConnection() (*gorm.DB, error) {
	db, err := openMySqlConnection(config.MysqlDSN())
	if err != nil {
		logrus.Fatalf("failed connect to database mysql : %s", err.Error())
	}

	return db, nil
}

func openMySqlConnection(dsn string) (*gorm.DB, error) {
	logrus.Infof("database MYSQL DSN : [%s]", dsn)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	logrus.Infof("success connected to database MYSQL [%s:%d]", config.MysqlHost(), config.MysqlPort())

	return db, nil
}
