package service

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("resource conflict")
	ErrUnavailable  = errors.New("resource unavailable")
)
