package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/utils"
	"sync"
)

type studentService struct {
	studentRepository entity.StudentRepository
	taskQueue         entity.TaskQueue
}

// NewStudentService is function to create new instance studentService
func NewStudentService(studentRepository entity.StudentRepository, taskQueue entity.TaskQueue) entity.StudentService {
	return &studentService{
		studentRepository: studentRepository,
		taskQueue:         taskQueue,
	}
}

func (s *studentService) InsertStudent(ctx context.Context, request *entity.RequestInsertStudent) error {
	//TODO implement me
	panic("implement me")
}

// InsertStudentBulk is method to insert student bulk
func (s *studentService) InsertStudentBulk(ctx context.Context, request *entity.RequestInsertStudentBulk) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
	})

	// validate request
	if err := request.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	var (
		wg = &sync.WaitGroup{}
		mu = &sync.Mutex{}
	)

	// looping each student
	for _, student := range request.Data {
		wg.Add(1)
		go func(wg *sync.WaitGroup, mu *sync.Mutex) {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()

			// create earch task for each student
			if err := s.taskQueue.Enqueue(ctx, entity.TaskInsertStudent, student); err != nil {
				logger.Error(err)
				return
			}
		}(wg, mu)
	}

	// wait all goroutines done
	wg.Wait()

	return nil
}
