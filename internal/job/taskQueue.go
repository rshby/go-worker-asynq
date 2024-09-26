package job

import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/config"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/utils"
)

type taskQueue struct {
	client    *asynq.Client
	namespace string
}

// NewTaskQueue is function to create new instance task queue
func NewTaskQueue(redisOpt asynq.RedisConnOpt) entity.TaskQueue {
	client := asynq.NewClient(redisOpt)

	return &taskQueue{
		client:    client,
		namespace: config.WorkerNamespace(),
	}
}

// Enqueue is method to enqueue task
func (t *taskQueue) Enqueue(ctx context.Context, taskName string, data any) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"data":    utils.Dump(data),
	})

	dataMarshalled, err := utils.JSONMarshal(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	// create new task
	task := asynq.NewTask(taskName, dataMarshalled)
	if _, err := t.client.Enqueue(task, t.generalOpts()...); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// Stop is method to stop taskQueue client
func (t *taskQueue) Stop() {
	if t.client != nil {
		if err := t.client.Close(); err != nil {
			logrus.Error(err)
		}

		logrus.Info("success stop TaskQueue 🔴")
	}
}

// generalOpts is method to get task options
func (t *taskQueue) generalOpts() []asynq.Option {
	return []asynq.Option{
		asynq.Queue(t.namespace),
		asynq.Retention(config.WorkerTaskRetention()),
		asynq.MaxRetry(config.WorkerRetryAttemps()),
		asynq.Timeout(config.WorkerTimeout()),
	}
}
