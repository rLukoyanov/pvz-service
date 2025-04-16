package errors

import "errors"

var (
	ErrCityNotAllowed = errors.New("city not allowed")
	ErrNotFound       = errors.New("not found")
	ErrInvalidInput   = errors.New("invalid input")
)
