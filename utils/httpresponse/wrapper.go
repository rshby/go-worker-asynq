package httpresponse

import (
	"github.com/gin-gonic/gin"
)

type WrapperDto struct {
	Message   string
	Data      any
	ErrorCode string
}

type HttpResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewHttpResponse() *HttpResponse {
	return new(HttpResponse)
}

func (h *HttpResponse) WithMessage(message string) *HttpResponse {
	h.Message = message
	return h
}

func (h *HttpResponse) WithData(data any) *HttpResponse {
	h.Data = data
	return h
}

type WrapperResponseDTO struct {
	HttpResponse
}

func (h *HttpResponse) ToWrapperResponseDTO(ctx *gin.Context, httpStatus int) {
	var wrapperResponseDTO WrapperResponseDTO
	wrapperResponseDTO.Data = h.Data
	wrapperResponseDTO.Message = h.Message

	ctx.JSON(httpStatus, wrapperResponseDTO)
}
