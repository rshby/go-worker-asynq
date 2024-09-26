package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// GetEnv is function to load env
func GetEnv(key string) string {
	// load env with godotenv
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal(err)
	}

	// get env from given key
	return os.Getenv(key)
}

// AppPort is function to get app port from env
func AppPort() int {
	if port := GetEnv("APP_PORT"); port != "" {
		appPort, err := strconv.Atoi(port)
		if err != nil {
			logrus.Error(err)
			return DefaultAppPort
		}

		return appPort
	}

	return DefaultAppPort
}

// MysqlUser is function to get mysql user from env
func MysqlUser() string {
	if user := GetEnv("MYSQL_USER"); user != "" {
		return user
	}

	return DefaultMysqlUser
}

// MysqlPassword is function to get mysql password from env
func MysqlPassword() string {
	if password := GetEnv("MYSQL_PASSWORD"); password != "" {
		return password
	}

	return DefaultMysqlPassword
}

// MysqlHost is function to get mysql host from env
func MysqlHost() string {
	if host := GetEnv("MYSQL_HOST"); host != "" {
		return host
	}

	return DefaultMysqlHost
}

// MysqlPort is function to get mysql port from env
func MysqlPort() int {
	if port := GetEnv("MYSQL_PORT"); port != "" {
		mysqlPort, err := strconv.Atoi(port)
		if err != nil {
			logrus.Error(err)
			return DefaultMysqlPort
		}

		return mysqlPort
	}

	return DefaultMysqlPort
}

// MysqlDbName is function to get mysql db name from env
func MysqlDbName() string {
	if dbName := os.Getenv("MYSQL_DB_NAME"); dbName != "" {
		return dbName
	}

	return DefaultMysqlDbName
}

// MysqlTimezone is function to get mysql timezone from env
func MysqlTimezone() string {
	if tz := GetEnv("MYSQL_TIMEZONE"); tz != "" {
		return tz
	}

	return DefaultMysqlTimezone
}

// MysqlDSN is function to get mysql dsn
func MysqlDSN() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		MysqlUser(), MysqlPassword(), MysqlHost(), MysqlPort(), MysqlDbName(), MysqlTimezone())

	return dsn
}
