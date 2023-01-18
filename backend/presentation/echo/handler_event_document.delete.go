package echo

import (
	"fmt"
	"net/http"
	"prc_hub_back/application/event"
	"prc_hub_back/domain/model/jwt"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// (DELETE /events/{id}/documents/{document_id})
func (*Server) DeleteEventsIdDocumentsDocumentId(ctx echo.Context) error {
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
	// Bind document_id
	var documentId DocumentId
	err = runtime.BindStyledParameterWithLocation("simple", false, "document_id", runtime.ParamLocationPath, ctx.Param("document_id"), &documentId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter document_id: %s", err))
	}

	// Delete document
	err = event.DeleteDocument(documentId, jcc.Id)
	if err != nil {
		return JSONMessage(ctx, event.ErrToCode(err), err.Error())
	}

	return JSONPretty(ctx, http.StatusNoContent, "success")
}
