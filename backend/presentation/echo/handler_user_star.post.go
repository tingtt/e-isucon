package echo

import (
	"fmt"
	"net/http"
	"prc_hub_back/application/user"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// (POST /users/{id}/star)
func (*Server) PostUsersIdStar(ctx echo.Context) error {
	// Bind id
	var id Id
	err := runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	count, err := user.AddStar(uint64(id))
	if err != nil {
		return JSONMessage(ctx, http.StatusInternalServerError, err.Error())
	}
	return JSONPretty(ctx, http.StatusOK, map[string]uint64{"count": count})
}
