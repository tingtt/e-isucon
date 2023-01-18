package echo

import (
	"errors"
	"net/http"
	"prc_hub_back/application/user"

	"github.com/labstack/echo/v4"
)

func (b LoginBody) Validate() error {
	if b.Email == "" {
		return errors.New("email required")
	}
	if b.Password == "" {
		return errors.New("password requied")
	}
	return nil
}

// (POST /users/sign_in)
func (*Server) PostUsersSignIn(ctx echo.Context) error {
	// Bind body
	body := new(PostUsersSignInJSONRequestBody)
	if err := ctx.Bind(body); err != nil {
		return JSONMessage(ctx, http.StatusBadRequest, err.Error())
	}

	// Verify and generate jwt
	token, verify, err := user.Verify(string(body.Email), body.Password)
	if err != nil {
		return JSONMessage(ctx, user.ErrToCode(err), err.Error())
	}
	if !verify {
		return JSONMessage(ctx, http.StatusUnauthorized, "failed to sign in")
	}
	return JSONPretty(ctx, http.StatusOK, Token{Token: token})
}
