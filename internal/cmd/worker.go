package cmd

import (
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-worker-asynq/cacher"
	"go-worker-asynq/config"
	"go-worker-asynq/internal/database"
	"go-worker-asynq/internal/job"
	"go-worker-asynq/utils"
	"os"
	"os/signal"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "run worker",
	Long:  "this command is to start worker",
	Run:   runWorker,
}

func init() {
	RootCmd.AddCommand(workerCmd)
}

// runWorker is function to run worker
func runWorker(cmd *cobra.Command, args []string) {
	// connect to database mysql
	db, err := database.InitializeMySqlConnection()
	if err != nil {
		logrus.Fatal(err)
	}

	dbMysql, err := db.DB()
	if err != nil {
		logrus.Fatal(err)
	}

	// create instance cacheManager
	cacheManager := cacher.ConstructCacheManager()
	cacheManager.SetDisableCaching(!config.EnableCache())
	if config.EnableCache() {
		// connect database redis cache
		redisDb, err := database.InitializeRedigoRedisConnectionPool(config.RedisCacheDSN(), nil)
		if err != nil {
			logrus.Fatalf("failed to connect REDIS [%s] : %s", config.RedisCacheDSN(), err.Error())
		}
		defer utils.WrapCloser(redisDb.Close)
		cacheManager.SetConnectionPool(redisDb)
	}

	// parse redis dsn
	redisOpt, err := asynq.ParseRedisURI(config.RedisWorkerDSN())
	if err != nil {
		logrus.Fatal(err)
	}

	if _, ok := redisOpt.(asynq.RedisConnOpt); !ok {
		logrus.Fatalf("failed to parse REDIS worker DSN [%s]", config.RedisWorkerDSN())
	}

	// create instance taskQueue
	taskQueue := job.NewTaskQueue(redisOpt)

	// create instance taskHandler
	taskHandler := job.NewTaskHandler(db, taskQueue)

	// create instance taskProcessor
	taskProcessor := job.NewTaskProcessor(redisOpt, config.WorkerNamespace(), taskHandler)

	// create channel signal interrupt
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go taskProcessor.Run()

	// waiting channel until receive signal interrupt
	<-sigCh
	logrus.Info("receive interrupt signal")

	close(sigCh)

	// stop taskQueue
	taskQueue.Stop()

	// stop taskProcessor
	taskProcessor.Stop()

	// close redis, use inside function gracefullShutdown
	gracefullShutdown(nil)

	// close mysql
	_ = dbMysql.Close()

	logrus.Info("worker exiting!")
}
