package echo

import (
	"fmt"
	"prc_hub_back/domain/model/jwt"
	"prc_hub_back/domain/model/logger"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Start(port uint, jwtIssuer string, jwtSecret string, allowOrigins []string) {
	// echoサーバーのインスタンス生成
	e := echo.New()

	// CORS
	if allowOrigins != nil {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: allowOrigins,
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
				echo.HeaderAuthorization,
			},
		}))
		logger.Logger().Info("cors enabled")
		logger.Logger().Debugf("cors allow origins: %v", allowOrigins)
	}

	// JWT
	jwt.InitWithSkipper(
		e,
		jwtIssuer,
		jwtSecret,
		func(c echo.Context) bool {
			// 公開エンドポイントのJWT認証をスキップ
			return c.Path() == "/users" && c.Request().Method == "POST" ||
				c.Path() == "/users/oauth2/:oauth_providers/register" && c.Request().Method == "POST" ||
				c.Path() == "/users/sign_in" && c.Request().Method == "POST" ||
				c.Path() == "/events" && c.Request().Method == "GET" ||
				c.Path() == "/events/:id" && c.Request().Method == "GET" ||
				c.Path() == "/events/:id/documents" && c.Request().Method == "GET" ||
				c.Path() == "/events/:id/documents/:document_id" && c.Request().Method == "GET" ||
				c.Path() == "/users/:id/star" && c.Request().Method == "POST" ||
				c.Path() == "/reset" && c.Request().Method == "POST"
		},
	)

	// handlerの登録
	var server *Server

	// ↓ スコア測定に直接関係するエンドポイント
	e.GET("/events", server.GetEvents)
	e.POST("/events", server.PostEvents)
	e.GET("/events/:id", server.GetEventsId)
	e.GET("/events/:id/documents", server.GetEventsIdDocuments)
	e.POST("/events/:id/documents", server.PostEventsIdDocuments)
	e.GET("/events/:id/documents/:document_id", server.GetEventsIdDocumentsDocumentId)
	e.POST("/reset", server.PostReset)
	e.GET("/users", server.GetUsers)
	e.POST("/users/sign_in", server.PostUsersSignIn)
	e.GET("/users/:id", server.GetUsersId)
	e.POST("/users/:id/star", server.PostUsersIdStar)

	// ↓ スコアに直接関係しないため他部分の変更による影響がない場合は原則変更しなくて構わない
	e.DELETE("/events/:id", server.DeleteEventsId)
	e.PATCH("/events/:id", server.PatchEventsId)
	e.DELETE("/events/:id/documents/:document_id", server.DeleteEventsIdDocumentsDocumentId)
	e.PATCH("/events/:id/documents/:document_id", server.PatchEventsIdDocumentsDocumentId)
	e.DELETE("/users", server.DeleteUsers)
	e.POST("/users", server.PostUsers)
	e.DELETE("/users/:id", server.DeleteUsersId)
	e.PATCH("/users/:id", server.PatchUsersId)

	// echoサーバーの起動
	logger.Logger().Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
