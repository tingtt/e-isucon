package jwt

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TODO: コンストラクタを整理
// シングルトン変数
var (
	jwtIssuer *string
	jwtSecret *string
)

// JWTミドルウェアの初期化
func Init(e *echo.Echo, issuer string, secret string) {
	jwtIssuer = &issuer
	jwtSecret = &secret
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwtCustumClaims{},
		SigningKey: []byte(secret),
	}))
}

// JWTミドルウェアの初期化
func InitWithSkipper(e *echo.Echo, issuer string, secret string, skipper func(c echo.Context) bool) {
	jwtIssuer = &issuer
	jwtSecret = &secret
	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwtCustumClaims{},
		SigningKey: []byte(secret),
		Skipper:    skipper,
	}))
}
