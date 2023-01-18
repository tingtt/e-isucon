package echo

import (
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/labstack/echo/v4"
)

// (POST /events)
func (*Server) PostEvents(ctx echo.Context) error {
	// Get jwt claim
	jcc, err := jwt.CheckProvided(ctx)
	if err != nil {
		return JSONMessage(ctx, http.StatusUnauthorized, err.Error())
	}

	// bind body
	body := new(event.CreateEventParam)
	if err := ctx.Bind(body); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}

	// Create event
	ed, err := event.CreateEvent(*body, jcc.Id)
	if err != nil {
		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusCreated, ed)
}
