package user

import "errors"

// Errors
var (
	ErrValidateNameCannotBeEmpty     = errors.New("user name cannot be empty")
	ErrValidateEmailCannotBeEmpty    = errors.New("user email cannot be empty")
	ErrValidateEmailAlreadyUsed      = errors.New("email already used")
	ErrValidatePasswordCannotBeEmpty = errors.New("user password cannot be empty")
	ErrValidatePasswordLength        = errors.New("password must be at least 8 characters")
)

func validateName(name string) error {
	// 空文字チェック
	if name == "" {
		return ErrValidateNameCannotBeEmpty
	}
	return nil
}

func validateEmail(email string) error {
	// 空文字チェック
	if email == "" {
		return ErrValidateEmailCannotBeEmpty
	}
	// 重複チェック

	_, err := GetByEmail(email)
	if err == nil || err != ErrUserNotFound {
		return ErrValidateEmailAlreadyUsed
	}
	if err == ErrUserNotFound {
		err = nil
	}
	return err
}

func validatePassword(password string) error {
	// 空文字チェック
	if password == "" {
		return ErrValidatePasswordCannotBeEmpty
	}
	// 文字長チェック
	if len(password) < 8 {
		return ErrValidatePasswordLength
	}
	return nil
}
