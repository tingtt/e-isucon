package echo

import "github.com/labstack/echo/v4"

type Server struct{}

var indent = "	"

func JSONMessage(ctx echo.Context, code int, msg string) error {
	return ctx.JSONPretty(code, map[string]string{"message": msg}, indent)
}

func JSONPretty(ctx echo.Context, code int, i interface{}) error {
	return ctx.JSONPretty(code, i, indent)
}
