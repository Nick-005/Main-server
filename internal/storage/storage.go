package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url was not found")
	ErrURLExists   = errors.New("url exists")
)
