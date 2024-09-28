package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/utils"
	"go-worker-asynq/utils/httpresponse"
	"net/http"
)

// InitStudentRoutes is method to register student routes
func (r *Router) InitStudentRoutes(app *gin.RouterGroup) {
	studentGroup := app.Group("student")
	{
		studentGroup.POST("/bulk", r.InsertStudentBulk)
	}
}

// InsertStudentBulk godoc
// Endpoint for insert bulk students
// @Summary Endpoint for insert bulk students
// @Description
// @Tags Student
// @Accept json
// @Produce json
// @Param Accept       header string true "Example : application/json"
// @Param Content-Type header string true "Example : application/json"
// @Param request body entity.RequestInsertStudentBulk true "payload request body"
// @Success 200 {object} entity.SwaggerResponseOKDTO{}
// @Router /v1/student/bulk [POST]
func (r *Router) InsertStudentBulk(c *gin.Context) {
	logger := logrus.WithContext(c).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(c),
	})

	// decode requestBody to object
	var request entity.RequestInsertStudentBulk
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(err)
		httpErrorHandler(c, err)
		return
	}

	// call method in service
	if err := r.studentService.InsertStudentBulk(c, &request); err != nil {
		logger.Error(err)
		httpErrorHandler(c, err)
		return
	}

	// success
	httpresponse.NewHttpResponse().WithMessage(successResponse["InsertStudentBulk"]).ToWrapperResponseDTO(c, http.StatusOK)
}
