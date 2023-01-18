package event

import "errors"

// Errors
var (
	ErrEventNotFound         = errors.New("event not found")
	ErrEventDocumentNotFound = errors.New("event document not found")
	ErrNoUpdates             = errors.New("no updates")
)
