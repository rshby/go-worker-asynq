package service

import "errors"

var (
	ErrNotFound            = errors.New("error not found")
	ErrBadRequest          = errors.New("error bad request")
	ErrInternalServerError = errors.New("error internal server")
	ErrStudentAlreadyExist = errors.New("student already exists in database")
)
