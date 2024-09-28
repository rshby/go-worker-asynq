package job

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/internal/job/service/interfaces"
	"go-worker-asynq/utils"
	"gorm.io/gorm"
)

type TaskHandler struct {
	db                *gorm.DB
	taskQueue         entity.TaskQueue
	jobStudentService interfaces.JobStudentService
}

// NewTaskHandler is function to create new instance taskHandler
func NewTaskHandler(db *gorm.DB, taskQueue entity.TaskQueue, jobStudentService interfaces.JobStudentService) *TaskHandler {
	return &TaskHandler{
		db:                db,
		taskQueue:         taskQueue,
		jobStudentService: jobStudentService,
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

	// call method in service to insert student
	if err := t.jobStudentService.InsertStudent(ctx, &student); err != nil {
		return err
	}

	return nil
}
