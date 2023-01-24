package echo

import "github.com/labstack/echo/v4"

type Server struct{}

func JSONMessage(ctx echo.Context, code int, msg string) error {
	return ctx.JSON(code, map[string]string{"message": msg})
}

func JSONPretty(ctx echo.Context, code int, i interface{}) error {
	return ctx.JSON(code, i)
}
