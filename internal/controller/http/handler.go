package http

import (
	"github.com/gin-gonic/gin"
	"go-worker-asynq/internal/entity"
)

type Router struct {
	studentService entity.StudentService
}

// RouteService is function to create router and register endpoints
func RouteService(app *gin.RouterGroup, studentService entity.StudentService) {
	router := Router{
		studentService: studentService,
	}

	// register endpoints
	router.Handlers(app)
}

// Handlers is method to register endpoints
func (r *Router) Handlers(app *gin.RouterGroup) {
	apiV1Group := app.Group("v1")
	{
		r.InitStudentRoutes(apiV1Group)
	}
}
