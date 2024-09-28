package cmd

import (
	"go-worker-asynq/cacher"
	"go-worker-asynq/internal/entity"
	jobService "go-worker-asynq/internal/job/service"
	jobServiceInterfaces "go-worker-asynq/internal/job/service/interfaces"
	"go-worker-asynq/internal/repository"
	"go-worker-asynq/internal/service"
	"gorm.io/gorm"
)

// InitStudentService is method to create studentService with dependency injection
func InitStudentService(db *gorm.DB, cache cacher.CacheManager, taskQueue entity.TaskQueue) entity.StudentService {
	// register repository
	studentRepository := repository.NewStudentRepository(db, cache)

	// create service
	return service.NewStudentService(studentRepository, taskQueue)
}

// InitJobStudentService is function to create instance jobStudentService with dependency injection
func InitJobStudentService(db *gorm.DB, cache cacher.CacheManager) jobServiceInterfaces.JobStudentService {
	// register repository
	studentRepository := repository.NewStudentRepository(db, cache)

	return jobService.NewJobStudentService(studentRepository)
}
