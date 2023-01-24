package user

import (
	"errors"
	"prc_hub_back/domain/model/user"
)

// Errors
var (
	ErrRepositoryNotInitialized = errors.New("repository not initialized")
)

func Init(u string, password string, db string) {
	user.InitRepository(u, password, db)
}
