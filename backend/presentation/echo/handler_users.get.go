package echo

import (
	"net/http"
	"prc_hub_back/application/user"
	"prc_hub_back/domain/model/jwt"

	"github.com/labstack/echo/v4"
)

// (GET /users)
func (*Server) GetUsers(ctx echo.Context) error {
	// Get jwt claim
	_, err := jwt.CheckProvided(ctx)
	if err != nil {
		return JSONMessage(ctx, http.StatusUnauthorized, err.Error())
	}

	// Bind query
	query := new(user.GetUserListQuery)
	if err := ctx.Bind(query); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}

	// Get users
	u, err := user.GetList(*query)
	if err != nil {
		return JSONMessage(ctx, user.ErrToCode(err), err.Error())
	}
	return JSONPretty(ctx, http.StatusOK, u)
}
