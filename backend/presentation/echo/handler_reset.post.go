package echo

import (
	"net/http"
	"prc_hub_back/application/eisucon"

	"github.com/labstack/echo/v4"
)

// (POST /reset)
func (*Server) PostReset(ctx echo.Context) error {
	err := eisucon.Migrate()
	if err != nil {
		return JSONMessage(ctx, http.StatusInternalServerError, err.Error())
	}
	return JSONMessage(ctx, http.StatusOK, "success")
}
