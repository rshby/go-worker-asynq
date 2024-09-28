package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/internal/job/service/interfaces"
	"go-worker-asynq/internal/service"
	"go-worker-asynq/utils"
)

type jobStudentService struct {
	studentRepository entity.StudentRepository
}

// NewJobStudentService is meth
func NewJobStudentService(studentRepository entity.StudentRepository) interfaces.JobStudentService {
	return &jobStudentService{
		studentRepository: studentRepository,
	}
}

// InsertStudent is method to insert student from taskProcessor
func (j *jobStudentService) InsertStudent(ctx context.Context, input *entity.Student) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(ctx),
		"input":   utils.Dump(input),
	})

	// lock proses insert
	mutexUnlock, err := j.studentRepository.LockInsertStudentByIdentityNumber(ctx, input.IdentityNumber)
	defer mutexUnlock()
	if err != nil {
		logger.Errorf("locked student %s : %s", input.IdentityNumber, err)
		return err
	}

	// get by identity_number, check if student already exist in database
	existStudent, err := j.studentRepository.GetStudentByIdentityNumber(ctx, input.IdentityNumber)
	if err != nil {
		logger.Error(err)
		return err
	}

	// if student with identity_number already exists in database
	if existStudent != nil {
		logger.Error(service.ErrStudentAlreadyExist)
		return service.ErrStudentAlreadyExist
	}

	// insert
	if _, err = j.studentRepository.InsertStudent(ctx, input); err != nil {
		logger.Error(err)
		return err
	}

	// success insert
	logrus.Infof("success insert data student with identity_number [%s]", input.IdentityNumber)
	return nil
}
