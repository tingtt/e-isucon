package user

import (
	"net/http"
	"prc_hub_back/domain/model/user"
)

var (
	// 404
	ErrUserNotFound = user.ErrUserNotFound

	// 405
	ErrAdminUserCannnotDelete = user.ErrAdminUserCannnotDelete

	// 400
	ErrValidateEmailAlreadyUsed        = user.ErrValidateEmailAlreadyUsed
	ErrPostEventAvailabledCannotUpdate = user.ErrPostEventAvailabledCannotUpdate
	ErrManageCannotUpdate              = user.ErrManageCannotUpdate
	ErrNoUpdates                       = user.ErrNoUpdates

	// 422
	ErrValidateNameCannotBeEmpty     = user.ErrValidateNameCannotBeEmpty
	ErrValidateEmailCannotBeEmpty    = user.ErrValidateEmailCannotBeEmpty
	ErrValidatePasswordCannotBeEmpty = user.ErrValidatePasswordCannotBeEmpty
	ErrValidatePasswordLength        = user.ErrValidatePasswordLength

	// 500
	ErrConflictUserStars = user.ErrConflictUserStars
)

func ErrToCode(e error) (code int) {
	switch e {
	case ErrUserNotFound:
		// 404
		return http.StatusNotFound
	case ErrAdminUserCannnotDelete:
		// 405
		return http.StatusMethodNotAllowed
	case ErrValidateEmailAlreadyUsed, ErrPostEventAvailabledCannotUpdate, ErrManageCannotUpdate, ErrNoUpdates:
		// 400
		return http.StatusBadRequest
	case ErrValidateNameCannotBeEmpty, ErrValidateEmailCannotBeEmpty, ErrValidatePasswordCannotBeEmpty, ErrValidatePasswordLength:
		// 422
		return http.StatusUnprocessableEntity
	default:
		// 500
		return http.StatusInternalServerError
	}
}
