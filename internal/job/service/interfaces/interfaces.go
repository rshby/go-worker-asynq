package interfaces

import (
	"context"
	"go-worker-asynq/internal/entity"
)

type JobStudentService interface {
	InsertStudent(ctx context.Context, input *entity.Student) error
}
