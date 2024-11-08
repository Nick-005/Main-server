package storage

import "errors"

var (
	ErrVACNotFound    = errors.New("vacancy was not found")
	ErrVACExists      = errors.New("vacancy already exists")
	ErrVACLimitIsOver = errors.New("limit for this org is over")
	ErrVACSomething   = errors.New("something went wrong with vacancy")
)
