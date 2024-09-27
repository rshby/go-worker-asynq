package cacher

import "time"

const (
	defaultTTL            = 15 * time.Minute
	defaultNilTTL         = 5 * time.Minute
	defaultLockDuration   = 1 * time.Hour
	defaultLockTries      = 1
	defaultWaitTime       = 15 * time.Second
	defaultPrefixCacheKey = "go-worker-asynq"
)
