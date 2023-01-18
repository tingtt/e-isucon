package user

import (
	"golang.org/x/crypto/bcrypt"
)

// パスワード検証
func (u *User) Verify(password string) (verify bool, err error) {
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}
