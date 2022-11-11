package utils

import "errors"

const (
	E_not_found      = "not_found"
	E_internal_error = "internal_error"
)

var NOT_FOUND_ERRORS error = errors.New("NOT_FOUND")
