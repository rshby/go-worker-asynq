package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
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

// EnvironmentMode is function to get env mode from env
func EnvironmentMode() string {
	return GetEnv("MODE")
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

// EnableCache is function to get enable cache from env
func EnableCache() bool {
	if enable := GetEnv("ENABLE_CACHE"); enable != "" {
		parseBool, err := strconv.ParseBool(enable)
		if err != nil {
			logrus.Error(err)
			return false
		}

		return parseBool
	}

	return false
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

// RedisCacheHost is function to get redis cache host from env
func RedisCacheHost() string {
	if host := GetEnv("REDIS_CACHE_HOST"); host != "" {
		return host
	}

	return DefaultRedisCacheHost
}

// RedisWorkerHost is function to get redis worker host from env
func RedisWorkerHost() string {
	if host := GetEnv("REDIS_WORKER_HOST"); host != "" {
		return host
	}

	return DefaultRedisWorkerHost
}

// RedisPort is function to get redis port from env
func RedisPort() int {
	if port := GetEnv("REDIS_PORT"); port != "" {
		redisPort, err := strconv.Atoi(port)
		if err != nil {
			logrus.Error(err)
			return DefaultRedisPort
		}

		return redisPort
	}

	return DefaultRedisPort
}

// RedisPingInterval is function to get redis ping interval from env
func RedisPingInterval() time.Duration {
	if ping := GetEnv("REDIS_PING_INTERVAL"); ping != "" {
		duration, err := time.ParseDuration(ping)
		if err != nil {
			logrus.Error(err)
			return DefaultRedisPingInterval
		}

		return duration
	}

	return DefaultRedisPingInterval
}

func RedisRetryAttemps() float64 {
	if retry := GetEnv("REDIS_RETRY_ATTEMPS"); retry != "" {
		redisRetry, err := strconv.ParseFloat(retry, 64)
		if err != nil {
			logrus.Error(err)
			return DefaultRedisRetryAttemps
		}

		return redisRetry
	}

	return DefaultRedisRetryAttemps
}

// RedisCacheDbNumber is function to get redis cache db number from env
func RedisCacheDbNumber() int {
	if dbNumber := GetEnv("REDIS_CACHE_DB_NUMBER"); dbNumber != "" {
		num, err := strconv.Atoi(dbNumber)
		if err != nil {
			logrus.Error(err)
			return DefaultRedisCacheDbNumber
		}

		return num
	}

	return DefaultRedisCacheDbNumber
}

// RedisWorkerDbNumber is function to get redis worker db number from env
func RedisWorkerDbNumber() int {
	if dbNumber := GetEnv("REDIS_WORKER_DB_NUMBER"); dbNumber != "" {
		num, err := strconv.Atoi(dbNumber)
		if err != nil {
			logrus.Error(err)
			return DefaultRedisWorkerDbNumber
		}

		return num
	}

	return DefaultRedisWorkerDbNumber
}

// RedisCacheDSN is function to get redis cache dsn
func RedisCacheDSN() string {
	dsn := fmt.Sprintf("redis://%s:%d/%d",
		RedisCacheHost(), RedisPort(), RedisCacheDbNumber())

	return dsn
}

// RedisWorkerDSN is function to get redis worker dsn
func RedisWorkerDSN() string {
	dsn := fmt.Sprintf("redis://%s:%d/%d",
		RedisWorkerHost(), RedisPort(), RedisWorkerDbNumber())

	return dsn
}

// WorkerNamespace is function to get worker namespace from env
func WorkerNamespace() string {
	return GetEnv("WORKER_NAMESPACE")
}

// WorkerTaskRetention is function to get worker task retention from env
func WorkerTaskRetention() time.Duration {
	if retention := GetEnv("WORKER_TASK_RETENTION"); retention != "" {
		duration, err := time.ParseDuration(retention)
		if err != nil {
			logrus.Error(err)
			return DefaultWorkerTaskRetention
		}

		return duration
	}

	return DefaultWorkerTaskRetention
}

// WorkerRetryAttemps is function to get worker retry attemps from env
func WorkerRetryAttemps() int {
	if retry := GetEnv("WORKER_RETRY_ATTEMPS"); retry != "" {
		workerRetry, err := strconv.Atoi(retry)
		if err != nil {
			logrus.Error(err)
			return DefaultWorkerRetryAttemps
		}

		return workerRetry
	}

	return DefaultWorkerRetryAttemps
}

// WorkerTimeout is function to get worker timeout from env
func WorkerTimeout() time.Duration {
	if timeout := GetEnv("WORKER_TIMEOUT"); timeout != "" {
		duration, err := time.ParseDuration(timeout)
		if err != nil {
			logrus.Error(err)
			return DefaultWorkerTimeout
		}

		return duration
	}

	return DefaultWorkerTimeout
}
