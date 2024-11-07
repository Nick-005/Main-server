package storage

import "errors"

var (
	ErrVACNotFound = errors.New("vacancy was not found")
	ErrVACExists   = errors.New("vacancy already exists")
)
