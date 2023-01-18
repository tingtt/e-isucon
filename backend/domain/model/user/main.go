package user

import "errors"

// Errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrNoUpdates         = errors.New("no updates")
	ErrConflictUserStars = errors.New("invalid result from query to count user stars")
)
