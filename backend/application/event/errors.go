package event

import (
	"net/http"
	"prc_hub_back/domain/model/event"
)

var (
	// 404
	ErrEventNotFound         = event.ErrEventNotFound
	ErrEventDocumentNotFound = event.ErrEventDocumentNotFound

	// 403
	ErrCannotCreateEvent         = event.ErrCannotCreateEvent
	ErrCannotUpdateEvent         = event.ErrCannotUpdateEvent
	ErrCannotDeleteEvent         = event.ErrCannotDeleteEvent
	ErrCannotCreateEventDocument = event.ErrCannotCreateEventDocument
	ErrCannotUpdateEventDocument = event.ErrCannotUpdateEventDocument
	ErrCannotDeleteEventDocument = event.ErrCannotDeleteEventDocument

	// 422
	ErrValidateEventTitleCannotBeEmpty         = event.ErrValidateEventTitleCannotBeEmpty
	ErrValidateDocumentNameCannotBeEmpty       = event.ErrValidateDocumentNameCannotBeEmpty
	ErrValidateUrlCannotBeEmpty                = event.ErrValidateUrlCannotBeEmpty
	ErrValidateEventDatetimesCannotBeEmpty     = event.ErrValidateEventDatetimesCannotBeEmpty
	ErrValidateEventDatetimeStartMustBeforeEnd = event.ErrValidateEventDatetimeStartMustBeforeEnd

	// 400
	ErrNoUpdates = event.ErrNoUpdates
)

func ErrToCode(e error) (code int) {
	switch e {
	case ErrEventNotFound,
		ErrEventDocumentNotFound:
		// 404
		return http.StatusNotFound
	case ErrCannotCreateEvent,
		ErrCannotUpdateEvent,
		ErrCannotDeleteEvent,
		ErrCannotCreateEventDocument,
		ErrCannotUpdateEventDocument,
		ErrCannotDeleteEventDocument:
		// 403
		return http.StatusForbidden
	case ErrValidateEventTitleCannotBeEmpty,
		ErrValidateDocumentNameCannotBeEmpty,
		ErrValidateUrlCannotBeEmpty,
		ErrValidateEventDatetimesCannotBeEmpty,
		ErrValidateEventDatetimeStartMustBeforeEnd:
		// 422
		return http.StatusUnprocessableEntity
	case ErrNoUpdates:
		// 400
		return http.StatusBadRequest
	default:
		// 500
		return http.StatusInternalServerError
	}
}
