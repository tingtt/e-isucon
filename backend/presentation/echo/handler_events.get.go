package echo

import (
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/labstack/echo/v4"
)

// (GET /events)
func (*Server) GetEvents(ctx echo.Context) error {
	// Get jwt claim
	var jwtId *int64
	jcc, err := jwt.Check(ctx)
	if err == nil {
		jwtId = &jcc.Id
	}

	// Bind query
	query := new(event.GetEventListQueryParam)
	type Query struct {
		Published       *bool   `query:"published"`
		Name            *string `query:"name"`
		NameContain     *string `query:"name_contain"`
		Location        *string `query:"location"`
		LocationContain *string `query:"location_contain"`
	}
	queryTmp := new(Query)
	if err := ctx.Bind(queryTmp); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}
	v := ctx.QueryParams()
	embed := v["embed"]
	query.Embed = &embed
	query.Name = queryTmp.Name
	query.NameContain = queryTmp.NameContain
	query.Location = queryTmp.Location
	query.LocationContain = queryTmp.LocationContain

	// Get events
	events, err := event.GetEventList(*query, jwtId)
	if err != nil {
		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
	}

	if events == nil {
		return JSONPretty(ctx, http.StatusOK, []interface{}{})
	}
	return JSONPretty(ctx, http.StatusOK, events)
}
