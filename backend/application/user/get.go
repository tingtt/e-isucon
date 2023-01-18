package user

import "prc_hub_back/domain/model/user"

func Get(id int64) (_ user.User, err error) {
	return user.Get(id)
}

func GetByEmail(email string) (user.User, error) {
	return user.GetByEmail(email)
}
