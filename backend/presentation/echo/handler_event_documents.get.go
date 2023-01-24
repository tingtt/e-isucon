package echo

import (
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/labstack/echo/v4"
)

// (GET /events/{id}/documents)
func (*Server) GetEventsIdDocuments(ctx echo.Context) error {
	// Get jwt claim
	jcc, err := jwt.Check(ctx)
	if err != nil {
		return JSONMessage(ctx, http.StatusUnauthorized, err.Error())
	}

	// Bind id
	id := ctx.Param("id")

	q := new(event.GetDocumentQueryParam)
	if err := ctx.Bind(q); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}

	// Get documents
	documents, err := event.GetDocumentList(event.GetDocumentQueryParam{
		EventId:     &id,
		Name:        q.Name,
		NameContain: q.NameContain,
	}, jcc.Id)
	if err != nil {
		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusOK, documents)
}
