package echo

import (
	"fmt"
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// (GET /events/{id})
func (*Server) GetEventsId(ctx echo.Context) error {
	// Get jwt claim
	var jwtId *int64
	jcc, err := jwt.Check(ctx)
	if err == nil {
		jwtId = &jcc.Id
	}

	// Bind id
	var id Id
	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Bind query
	v := ctx.QueryParams()
	embed := v["embed"]
	query := new(event.GetEventQueryParam)
	query.Embed = &embed

	// Get event
	e, err := event.GetEvent(id, *query, jwtId)
	if err != nil {
		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusOK, e)
}
