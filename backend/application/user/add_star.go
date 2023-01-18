package user

import "prc_hub_back/domain/model/user"

func AddStar(userId uint64) (count uint64, err error) {
	return user.AddStar(userId)
}
