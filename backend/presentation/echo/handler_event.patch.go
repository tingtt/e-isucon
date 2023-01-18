package echo

import (
	"fmt"
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// (PATCH /events/{id})
func (*Server) PatchEventsId(ctx echo.Context) error {
	// Get jwt claim
	jcc, err := jwt.CheckProvided(ctx)
	if err != nil {
		return JSONMessage(ctx, http.StatusUnauthorized, err.Error())
	}

	// Bind id
	var id Id
	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Bind body
	body := new(event.UpdateEventParam)
	if err := ctx.Bind(body); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}

	// Update event
	e, err := event.UpdateEvent(id, *body, jcc.Id)
	if err != nil {
		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusOK, e)
}
