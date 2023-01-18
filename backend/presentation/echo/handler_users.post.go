package echo

import (
	"net/http"
	"prc_hub_back/application/user"

	"github.com/labstack/echo/v4"
)

// (POST /users)
func (*Server) PostUsers(ctx echo.Context) error {
	// Bind body
	body := new(user.CreateUserParam)
	if err := ctx.Bind(body); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}

	// Create user
	uwt, err := user.Create(*body)
	if err != nil {
		return JSONMessage(ctx, user.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusCreated, uwt)
}
