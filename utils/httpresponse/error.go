package httpresponse

import (
	"context"
	"go-worker-asynq/utils"

	"github.com/sirupsen/logrus"
)

type WrapperErrorResponseDTO struct {
	WrapperResponseDTO
	ErrorCode string `json:"errorCode,omitempty"`
}

type HTTPError struct {
	Code      int
	Message   string
	ErrorCode string
	Data      any `json:"data,omitempty"`
}

func NewHTTPError() *HTTPError {
	return new(HTTPError)
}

func (httpError *HTTPError) Error() string {
	if httpError == nil {
		return ""
	}

	return httpError.Message
}

func (httpError *HTTPError) StatusCode() int {
	if httpError == nil {
		return 0
	}

	return httpError.Code
}

func (httpError *HTTPError) GetErrorCode() string {
	if httpError == nil {
		return ""
	}

	return httpError.ErrorCode
}

func (httpError *HTTPError) WithCode(httpStatusCode int) *HTTPError {
	httpError.Code = httpStatusCode
	return httpError
}

func (httpError *HTTPError) WithMessage(err error) *HTTPError {
	httpError.Message = err.Error()
	return httpError
}

func (httpError *HTTPError) ToResponseWithContext(context context.Context) WrapperErrorResponseDTO {
	logrus.WithContext(context).WithField("httpError", utils.Dump(httpError)).Error(httpError.Error())

	var wrapperErrorResponse WrapperErrorResponseDTO

	wrapperErrorResponse.Message = httpError.Message
	wrapperErrorResponse.ErrorCode = httpError.ErrorCode
	wrapperErrorResponse.Data = httpError.Data

	return wrapperErrorResponse
}

var (
	// 400 bad request
	ErrorBadRequest = HTTPError{ErrorCode: "CUST400000", Message: "Bad Request"}

	// 401 unauthorize
	ErrorInvalidSignature = HTTPError{ErrorCode: "CUST401000", Message: "Invalid Signature"}

	// 403 forbidden
	ErrorEmailTemplateInactive = HTTPError{ErrorCode: "CUST403000", Message: "Email Template is Inactive"}

	// 404 not found
	ErrSMTPConfigNotFound    = HTTPError{ErrorCode: "CUST404010", Message: "SMTP Config is Not Found"}
	ErrEmailPurposeNotFound  = HTTPError{ErrorCode: "CUST404020", Message: "Email Purpose is Not Found"}
	ErrServiceNotFound       = HTTPError{ErrorCode: "CUST404030", Message: "Service is Not Found"}
	ErrEmailTemplateNotFound = HTTPError{ErrorCode: "CUST404040", Message: "Email Template is Not Found"}

	// 500 internal server
	ErrorInternalServerError = HTTPError{ErrorCode: "CUST500000", Message: "Internal Server Error"}
)
