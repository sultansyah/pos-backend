package custom

import (
	"errors"
)

var (
	ErrAlreadyExists     = errors.New("resource already exists")
	ErrNotFound          = errors.New("resource not found")
	ErrInternal          = errors.New("internal server error")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("you are not authorized to access this resource")
	ErrImageRequired     = errors.New("image is required")
	ErrConflict          = errors.New("duplicate entry for key")
	ErrInsufficientStock = errors.New("insufficient stock")
)
