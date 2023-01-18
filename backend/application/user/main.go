package user

import (
	"errors"
	"prc_hub_back/domain/model/user"
)

// Errors
var (
	ErrRepositoryNotInitialized = errors.New("repository not initialized")
)

func Init(u string, password string, host string, port uint, db string) {
	user.InitRepository(u, password, host, port, db)
}
