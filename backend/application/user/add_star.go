package user

import "prc_hub_back/domain/model/user"

func AddStar(UserId string) (count uint64, err error) {
	return user.AddStar(UserId)
}
