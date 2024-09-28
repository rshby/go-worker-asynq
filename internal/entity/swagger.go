package entity

type SwaggerResponseOKDTO struct {
	Message string `json:"message" example:"success"`
	Data    any    `json:"data"`
}
