package job

import (
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"log"
)

type TaskProcessor struct {
	server      *asynq.Server
	taskHandler *TaskHandler
	namespace   string
}

// NewTaskProcessor is function to create new instance taskProcessor
func NewTaskProcessor(redisOpt asynq.RedisConnOpt, namespace string, tasHandler *TaskHandler) *TaskProcessor {
	logger := logrus.WithFields(logrus.Fields{
		"namespace": namespace,
	})

	// create server
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 100,
			Queues: map[string]int{
				namespace: 10,
			},
			Logger: logger,
		})

	// create instance taskProcessor
	taskProcessor := TaskProcessor{
		server:      server,
		taskHandler: tasHandler,
		namespace:   namespace,
	}

	return &taskProcessor
}

// Run is method to run task processor
func (t *TaskProcessor) Run() {
	t.RegisterTask()

	logrus.Info("Running Worker 🟢")
	if err := t.server.Start(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

// Stop is method to stop task processor
func (t *TaskProcessor) Stop() {
	if t.server != nil {
		t.server.Stop()
		t.server.Shutdown()
	}
}

// RegisterTask is method to register task and their handle function
func (t *TaskProcessor) RegisterTask() {
	mux.HandleFunc(entity.TaskInsertStudent, t.taskHandler.HandleTaskInsertStudent)
}
