package errors

import "errors"

var (
	ErrCityNotAllowed     = errors.New("недопустимый город")
	ErrCategoryNotAllowed = errors.New("недопустимая категория")
	ErrNotFound           = errors.New("не найдено")
	ErrInvalidInput       = errors.New("не верный ввод")
	ErrNoReceprionsFound  = errors.New("не нашли открытых приемок")
)
