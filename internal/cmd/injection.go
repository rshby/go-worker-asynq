package cmd

import (
	"go-worker-asynq/cacher"
	"go-worker-asynq/internal/entity"
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
