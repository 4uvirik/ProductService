package pkg

import (
	"errors"
	"net/http"
)

// Ошибка ответа
type ErrorResponse struct {
	Error string `json:"error"`
}

// MapErrorToStatus принимает ошибку из бизнес-/репозиторного слоя
// и возвращает соответствующий HTTP-статус:
// ErrNotFound → 404 Not Found
// остальные ошибки → 400 Bad Request
func MapErrorToStatus(err error) int {
	switch {
	case errors.Is(err, ErrNoFound):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}

var ErrNoFound = errors.New("no found")
