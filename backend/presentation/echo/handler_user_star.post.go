package echo

import (
	"net/http"
	"prc_hub_back/application/user"

	"github.com/labstack/echo/v4"
)

// (POST /users/{id}/star)
func (*Server) PostUsersIdStar(ctx echo.Context) error {
	// Bind id
	id := ctx.Param("id")

	count, err := user.AddStar(id)
	if err != nil {
		return JSONMessage(ctx, http.StatusInternalServerError, err.Error())
	}
	return JSONPretty(ctx, http.StatusOK, map[string]uint64{"count": count})
}
