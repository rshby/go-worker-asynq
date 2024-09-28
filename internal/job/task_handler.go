package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/utils"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db        *gorm.DB
	taskQueue entity.TaskQueue
}

// NewTaskHandler is function to create new instance taskHandler
func NewTaskHandler(db *gorm.DB, taskQueue entity.TaskQueue) *TaskHandler {
	return &TaskHandler{
		db:        db,
		taskQueue: taskQueue,
	}
}

// HandleTaskInsertStudent is method to handle task insert student data
func (t *TaskHandler) HandleTaskInsertStudent(ctx context.Context, task *asynq.Task) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"payload": string(task.Payload()),
	})

	// decode payload
	var student entity.Student
	if err := json.Unmarshal(task.Payload(), &student); err != nil {
		logger.Error(err)
		return err
	}

	fmt.Println(student)
	return nil
}
