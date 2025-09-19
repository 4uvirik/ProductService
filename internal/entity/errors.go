package entity

import "errors"

var ErrNotFound = errors.New("not found")
var ErrBadRequest = errors.New("bad request")
var ErrNoFields = errors.New("no fields")
var ErrInvalidConfig = errors.New("invalid config")
var ErrDBError = errors.New("db error")
