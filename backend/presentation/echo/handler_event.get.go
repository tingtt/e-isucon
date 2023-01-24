package echo

import (
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/labstack/echo/v4"
)

// (GET /events/{id})
func (*Server) GetEventsId(ctx echo.Context) error {
	// Get jwt claim
	var jwtId *string
	jcc, err := jwt.Check(ctx)
	if err == nil {
		jwtId = &jcc.Id
	}

	// Bind id
	id := ctx.Param("id")

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
