package exception

import (
	"errors"
	"fmt"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrBadRequest   = errors.New("bad request")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
	ErrNotFound     = errors.New("not found")
	ErrInternal     = errors.New("internal error")
)

type AppError struct {
	Kind    error
	Message string
	Err     error
}

func New(kind error, message string) error {
	return &AppError{
		Kind:    kind,
		Message: message,
	}
}

func Wrap(kind error, message string, err error) error {
	return &AppError{
		Kind:    kind,
		Message: message,
		Err:     err,
	}
}

func (err *AppError) Error() string {
	if err.Message != "" {
		return err.Message
	}
	if err.Err != nil {
		return err.Err.Error()
	}
	if err.Kind != nil {
		return err.Kind.Error()
	}
	return "application error"
}

func (err *AppError) Unwrap() error {
	return err.Err
}

func (err *AppError) Is(target error) bool {
	return target == err.Kind
}

func (err *AppError) FullError() string {
	if err.Err == nil {
		return err.Error()
	}
	return fmt.Sprintf("%s: %v", err.Error(), err.Err)
}
