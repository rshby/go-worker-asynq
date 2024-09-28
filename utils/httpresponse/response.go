package httpresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(ctx *gin.Context, error *HTTPError) {
	ctx.JSON(error.StatusCode(), error.ToResponseWithContext(ctx))
}

func OK(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusOK)
}

func Created(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusCreated)
}

func NoContent(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusNoContent)
}

func Accepted(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusAccepted)
}

func Unauthorized(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusUnauthorized)
}

func Forbidden(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusForbidden)
}

func NotFound(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusNotFound)
}

func BadRequest(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusBadRequest)
}

func InternalServerError(ctx *gin.Context, httpResponse *HttpResponse) {
	httpResponse.ToWrapperResponseDTO(ctx, http.StatusInternalServerError)
}
