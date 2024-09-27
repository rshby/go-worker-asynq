package database

import (
	"errors"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jpillora/backoff"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go-worker-asynq/config"
	"go-worker-asynq/utils"
	"time"
)

var (
	RedisConnPoll *redigo.Pool
	StopTickerCh  = make(chan bool)
)

// RedisConnectionPoolOptions options for the redis connection
type RedisConnectionPoolOptions struct {
	// Dial timeout for establishing new connections.
	// Default is 5 seconds. Only for go-redis.
	DialTimeout time.Duration

	// Enables read-only commands on slave nodes.
	ReadOnly bool

	// Timeout for socket reads. If reached, commands will fail
	// with a timeout instead of blocking. Use value -1 for no timeout and 0 for default.
	// Default is 3 seconds. Only for go-redis.
	ReadTimeout time.Duration

	// Timeout for socket writes. If reached, commands will fail
	// with a timeout instead of blocking.
	// Default is ReadTimeout. Only for go-redis.
	WriteTimeout time.Duration

	// Number of idle connections in the pool.
	IdleCount int

	// Maximum number of connections allocated by the pool at a given time.
	// When zero, there is no limit on the number of connections in the pool.
	PoolSize int

	// Close connections after remaining idle for this duration. If the value
	// is zero, then idle connections are not closed. Applications should set
	// the timeout to a value less than the server's timeout.
	IdleTimeout time.Duration

	// Close connections older than this duration. If the value is zero, then
	// the pool does not close connections based on age.
	MaxConnLifetime time.Duration
}

var defaultRedisConnectionPoolOptions = &RedisConnectionPoolOptions{
	IdleCount:       20,
	PoolSize:        100,
	IdleTimeout:     20 * time.Minute,
	MaxConnLifetime: 30 * time.Minute,
	DialTimeout:     10 * time.Second,
	WriteTimeout:    5 * time.Second,
	ReadTimeout:     5 * time.Second,
}

// InitializeRedigoRedisConnectionPool uses redigo library to establish the redis connection pool
func InitializeRedigoRedisConnectionPool(url string, opt *RedisConnectionPoolOptions) (*redigo.Pool, error) {
	if !isValidRedisStandaloneURL(url) {
		log.Fatal("invalid redis url :", url)
	}

	options := applyRedisConnectionPoolOptions(opt)

	redigoPool := redigo.Pool{
		MaxIdle:     options.IdleCount,
		MaxActive:   options.PoolSize,
		IdleTimeout: options.IdleTimeout,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.DialURL(url)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			return c, err
		},
		MaxConnLifetime: options.MaxConnLifetime,
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				log.Error(err)
			}
			return err
		},
		Wait: true, // wait for connection available when maxActive is reached
	}

	// test ping redis
	client := redigoPool.Get()
	reply, err := client.Do("PING")
	defer client.Close()

	if err != nil {
		return nil, err
	}

	if reply == nil {
		errMessage := utils.WriteStringTemplate("failed to connect REDIS [%s]", config.RedisCacheHost())
		log.Error(errMessage)
		return nil, errors.New(errMessage)
	}

	log.Infof("success connected to REDIS [%s]", config.RedisCacheDSN())
	RedisConnPoll = &redigoPool

	// continously check redis connection
	go checkRedisConnection(url, time.NewTicker(config.RedisPingInterval()))

	return RedisConnPoll, nil
}

// checkRedisConnection is function to continously check redis connection
func checkRedisConnection(url string, ticker *time.Ticker) {
	log.Info("start continously check REDIS ⏳")

	for {
		select {
		case <-StopTickerCh:
			ticker.Stop()
			close(StopTickerCh)
			log.Info("stop continously check REDIS 🔴")
			return
		case <-ticker.C:
			client := RedisConnPoll.Get()

			if _, err := client.Do("PING"); err != nil {
				reconnectRedisConnection(url)
			}

			if err := client.Close(); err != nil {
				log.Error(err)
			}
		}
	}
}

// reconnectRedisConnection is function to reconnect redis connection
func reconnectRedisConnection(url string) {
	b := backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    100 * time.Millisecond,
		Max:    1 * time.Second,
	}

	redisRetryAttempts := config.RedisRetryAttemps()

	for b.Attempt() < redisRetryAttempts {
		options := applyRedisConnectionPoolOptions(nil)

		client := &redigo.Pool{
			MaxIdle:     options.IdleCount,
			MaxActive:   options.PoolSize,
			IdleTimeout: options.IdleTimeout,
			Dial: func() (redigo.Conn, error) {
				c, err := redigo.DialURL(url)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			MaxConnLifetime: options.MaxConnLifetime,
			TestOnBorrow: func(c redigo.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Wait: true, // wait for connection available when maxActive is reached
		}

		if _, err := client.Get().Do("PING"); err == nil {
			// assign RedisConnPoll with new value redigo.Poll
			RedisConnPoll = client
			log.Warn("new redis connection")
			break
		}

		time.Sleep(b.Duration())
	}

	if b.Attempt() >= redisRetryAttempts {
		log.Fatal("maximum retry to connect database")
	}

	b.Reset()
}

func isValidRedisStandaloneURL(url string) bool {
	_, err := redis.ParseURL(url)
	if err != nil {
		log.Error(err)
	}

	return err == nil
}

func applyRedisConnectionPoolOptions(opt *RedisConnectionPoolOptions) *RedisConnectionPoolOptions {
	if opt != nil {
		return opt
	}

	return defaultRedisConnectionPoolOptions
}
