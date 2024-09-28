package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-worker-asynq/cacher"
	"go-worker-asynq/config"
	controllerHttp "go-worker-asynq/internal/controller/http"
	"go-worker-asynq/internal/database"
	"go-worker-asynq/internal/job"
	"go-worker-asynq/utils"
	"net/http"
	"os"
	"os/signal"
	"time"
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

// server is function to run server program
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
		logrus.Fatalf("failed to parse REDIS WORKER DSN [%s]", config.RedisWorkerDSN())
	}

	// create instance taskQueue
	taskQueue := job.NewTaskQueue(redisOpt)

	// register service
	studentService := InitStudentService(db, cacheManager, taskQueue)

	// create gin app
	app := gin.Default()
	app.UseRawPath = true
	app.UnescapePathValues = true
	app.RemoveExtraSlash = true

	// create router and register endpoints
	controllerHttp.RouteService(&app.RouterGroup, studentService)

	// create server
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", config.AppPort()),
		Handler: app,
	}

	// channel interrupt
	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	quitChan := make(chan bool, 1)
	signal.Notify(sigChan, os.Interrupt)

	// goroutine to check and gracefull shutdown server
	go func() {
		var err error
		for {
			select {
			case <-sigChan:
				logrus.Info("receive interrup signal ⚠️")
				taskQueue.Stop()
				gracefullShutdown(&srv)
				gracefullDbMYSQL(dbMysql)
				quitChan <- true
				return
			case err = <-errChan:
				logrus.Error(err)
				taskQueue.Stop()
				gracefullShutdown(&srv)
				gracefullDbMYSQL(dbMysql)
				quitChan <- true
				return
			}
		}
	}()

	// goroutine to run server
	go func() {
		var err error
		if err = srv.ListenAndServe(); err != nil {
			errChan <- err
			return
		}
	}()

	// wait quit channel
	<-quitChan

	// closing all channel
	close(sigChan)
	close(errChan)
	close(quitChan)

	logrus.Info("server exit ❌")
}

func gracefullShutdown(srv *http.Server) {
	// set value to chnannel, to stop redis check db
	database.StopTickerCh <- true
	time.Sleep(100 * time.Millisecond)

	if srv != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// shutdown server
		if err := srv.Shutdown(ctx); err != nil {
			logrus.Error(err)

			// if any error when shutdown, then force to close server
			if err = srv.Close(); err != nil {
				logrus.Fatalf("force close server 🔴")
			}
		}
		logrus.Info("success shutdown server 🔴")

		// close server after shutdown
		if err := srv.Close(); err != nil {
			logrus.Fatalf("force close server : %s", err.Error())
		}
		logrus.Info("success close server 🔴")
	}
}

func gracefullDbMYSQL(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			logrus.Fatalf("force close database : %s", err.Error())
		}

		logrus.Info("success stop db mysql 🔴")
	}
}
