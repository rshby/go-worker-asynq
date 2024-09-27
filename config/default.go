package config

import "time"

const (
	DefaultAppPort = 4000

	DefaultMysqlUser     = "root"
	DefaultMysqlPassword = "root"
	DefaultMysqlHost     = "localhost"
	DefaultMysqlPort     = 3306
	DefaultMysqlDbName   = "go_worker_asynq_db"
	DefaultMysqlTimezone = "UTC"

	DefaultRedisCacheHost      = "localhost"
	DefaultRedisWorkerHost     = "localhost"
	DefaultRedisPort           = 6379
	DefaultRedisPingInterval   = 5 * time.Second
	DefaultRedisRetryAttemps   = float64(3)
	DefaultRedisCacheDbNumber  = 0
	DefaultRedisWorkerDbNumber = 1

	DefaultWorkerTaskRetention = 2 * time.Hour
	DefaultWorkerRetryAttemps  = 3
	DefaultWorkerTimeout       = 2 * time.Hour
)
