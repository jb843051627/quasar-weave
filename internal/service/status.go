package service

import (
	"errors"
	"net/http"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func StatusCode(err error) int {
	switch {
	case err == model.ErrNotFound:
		return http.StatusNotFound
	case errors.Is(err, model.ErrInvalidState), errors.Is(err, model.ErrAlreadyExists), errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, model.ErrQueueClosed), errors.Is(err, ErrUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
