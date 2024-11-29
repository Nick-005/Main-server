package storage

import "errors"

var (
	ErrVACNotFound    = errors.New("vacancy was not found")
	ErrVACExists      = errors.New("vacancy already exists")
	ErrVACLimitIsOver = errors.New("limit for this org is over")
	ErrVACSomething   = errors.New("something went wrong with vacancy")

	ErrUSERNotFound    = errors.New("user was not found")
	ErrUSERExists      = errors.New("user already exists")
	ErrUSERLimitIsOver = errors.New("user for this org is over")
	ErrUSERSomething   = errors.New("something went wrong with user")
)
