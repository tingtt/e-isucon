package echo

import (
	"fmt"
	"net/http"
	"prc_hub_back/application/user"
	"prc_hub_back/domain/model/jwt"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// (GET /users/{id})
func (*Server) GetUsersId(ctx echo.Context) error {
	// Get jwt claim
	_, err := jwt.CheckProvided(ctx)
	// jcc, err := jwt.CheckProvided(ctx)
	if err != nil {
		return JSONMessage(ctx, http.StatusUnauthorized, err.Error())
	}

	// Bind id
	var id Id
	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Get user
	u, err := user.Get(id)
	if err != nil {
		return JSONMessage(ctx, user.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusOK, u)
}
