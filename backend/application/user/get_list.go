package user

import "prc_hub_back/domain/model/user"

type GetUserListQuery user.GetUserListQueryParam

func GetList(q GetUserListQuery) ([]user.User, error) {
	return user.GetList(user.GetUserListQueryParam(q))
}
