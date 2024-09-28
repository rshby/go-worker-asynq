package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go-worker-asynq/internal/service"
	"go-worker-asynq/utils"
	"go-worker-asynq/utils/httpresponse"
	"net/http"
)

var (
	ErrInternalServerError = httpresponse.NewHTTPError().WithCode(http.StatusInternalServerError).WithMessage(service.ErrInternalServerError)

	errorResponse = map[error]*httpresponse.HTTPError{
		service.ErrNotFound:            httpresponse.NewHTTPError().WithCode(http.StatusNotFound).WithMessage(service.ErrNotFound),
		service.ErrBadRequest:          httpresponse.NewHTTPError().WithCode(http.StatusBadRequest).WithMessage(service.ErrBadRequest),
		service.ErrInternalServerError: ErrInternalServerError,
	}

	successResponse = map[string]string{
		"InsertStudentBulk": "Success Process Insert Student Bulk",
	}
)

type (
	metaInfo struct {
		TotalItems int `json:"totalItems,omitempty"`
		TotalPages int `json:"totalPages,omitempty"`
	}

	paginationResponse[T any] struct {
		Items []T `json:"items"`
		metaInfo
	}
)

func httpErrorHandler(context *gin.Context, err error) {
	if err == nil {
		return
	}

	var validatorError validator.ValidationErrors

	switch {
	case errors.As(err, &validatorError):
		type validationErrorResponse struct {
			Field string
			Error string
		}

		validationErrors := err.(validator.ValidationErrors)
		errorsData := make([]validationErrorResponse, len(validationErrors))

		for i, fe := range validationErrors {
			errorsData[i] = validationErrorResponse{fe.Field(), msgForTag(fe)}
		}

		httpresponse.Error(context, &httpresponse.HTTPError{
			Code:    http.StatusBadRequest,
			Message: service.ErrBadRequest.Error(),
			Data:    errorsData,
		})

		return
	default:
		httpError, ok := errorResponse[err]
		if !ok {
			httpresponse.Error(context, ErrInternalServerError)
			return
		}

		httpresponse.Error(context, httpError)
		return
	}
}

func msgForTag(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "isUrl":
		return "Invalid URL format"
	case "min":
		return utils.WriteStringTemplate("The value must be greater than or equal to %s", fieldError.Param())
	case "max":
		return utils.WriteStringTemplate("The value must be less than or equal to %s", fieldError.Param())
	case "unique":
		return "The value must be unique"
	case "number":
		return "The value must be a number format"
	case "len":
		return utils.WriteStringTemplate("Length of the value must be %s character", fieldError.Param())

	}

	return utils.WriteStringTemplate("Validation failed for the '%s' tag", fieldError.Tag())
}
