package httperr

import (
	"errors"
	"github.com/4uvirik/ProductService/internal/entity"
	"net/http"
)

// ErrorResponse - ошибка ответа
type ErrorResponse struct {
	Error string `json:"error"`
}

// MapErrorToStatus принимает ошибку из бизнес или репозиторного слоя
// и возвращает соответствующий HTTP статус:
// ErrNotFound -> 404 Not Found
// остальные ошибки -> 400 Bad Request
func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, entity.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, entity.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, entity.ErrNoFields):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
