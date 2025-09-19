package httperr

import (
	"errors"
	"net/http"

	"github.com/4uvirik/ProductService/internal/entity"
)

// ErrorResponse - ошибка ответа.
type ErrorResponse struct {
	Error string `json:"error"`
}

// MapErrorToStatus принимает ошибку из бизнес или репозиторного слоя.
func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, entity.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, entity.ErrNoFields):
		return http.StatusBadRequest
	case errors.Is(err, entity.ErrBadRequest):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
