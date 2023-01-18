package jwt

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

// エラー
var (
	ErrTokenInvalid            = errors.New("token invalid")
	ErrTokenExpired            = errors.New("token expired")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrBearerTokenNotFound     = errors.New("bearer token not specified")
)

type jwtCustumClaims struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

type GenerateTokenParam struct {
	Id    int64
	Email string
	Admin bool
}

// トークン生成
func GenerateToken(p GenerateTokenParam) (token string, err error) {
	// claimの作成
	claims := &jwtCustumClaims{
		p.Id,
		p.Email,
		p.Admin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    *jwtIssuer,
		},
	}

	// トークンを生成
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return newToken.SignedString([]byte(*jwtSecret))
}

// トークンの検証
func verify(token *jwt.Token) (claims *jwtCustumClaims, err error) {
	claims = token.Claims.(*jwtCustumClaims)

	if !claims.VerifyIssuer(*jwtIssuer, true) {
		// 不正なトークン
		err = ErrTokenInvalid
		return
	}

	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		// 期限切れのトークン
		err = ErrTokenExpired
		return
	}

	return
}

// JWT middleware が有効なエンドポイントのトークン検証
func CheckProvided(ctx echo.Context) (claims *jwtCustumClaims, err error) {
	user := ctx.Get("user").(*jwt.Token)
	return verify(user)
}

// JWT middleware が無効なエンドポイントのトークン検証
func Check(ctx echo.Context) (claims *jwtCustumClaims, err error) {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		var token *jwt.Token
		token, err = jwt.ParseWithClaims(
			tokenString,
			&jwtCustumClaims{},
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, ErrUnexpectedSigningMethod
				}
				return []byte(*jwtSecret), nil
			},
		)
		if err != nil {
			return nil, err
		}
		return verify(token)
	}
	return nil, ErrBearerTokenNotFound
}
