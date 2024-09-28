package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-worker-asynq/internal/entity"
	"go-worker-asynq/utils"
)

// InitStudentRoutes is method to register student routes
func (r *Router) InitStudentRoutes(app *gin.RouterGroup) {
	studentGroup := app.Group("student")
	{
		studentGroup.POST("/bulk", r.InsertStudentBulk)
	}
}

func (r *Router) InsertStudentBulk(c *gin.Context) {
	logger := logrus.WithContext(c).WithFields(logrus.Fields{
		"context": utils.DumpIncomingContext(c),
	})

	// decode requestBody to object
	var request entity.RequestInsertStudentBulk
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error(err)
		// TODO : httpErrorHandler(c, err)
		return
	}
}
