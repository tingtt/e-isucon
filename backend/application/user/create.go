package user

import "prc_hub_back/domain/model/user"

type (
	CreateUserParam user.CreateUserParam
)

func Create(p CreateUserParam) (user.UserWithToken, error) {
	return user.CreateUser(
		user.CreateUserParam{
			Name:           p.Name,
			Email:          p.Email,
			Password:       p.Password,
			TwitterId:      p.TwitterId,
			GithubUsername: p.GithubUsername,
		},
	)
}
