package api_handler

import (
	"github.com/chur-squad/loveframe-server/auth"
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/labstack/echo/v4"
	"net/http"
)

type loginRequest struct {
	Id int64 `json:"id"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) RequestLogin(c echo.Context) error {
	ctx := c.(_context.EchoContext)

	request := new(loginRequest)
	if err := c.Bind(request); err != nil {
		return _error.WrapError(err)
	}

	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlAllowMethods:  "GET, POST, OPTIONS",
		echo.HeaderAccessControlRequestMethod: "*",
		echo.HeaderContentType:                "application/json",
	}

	user, err := h.parent.Mysql.UserById(ctx, request.Id)
	if err != nil {
		return _error.WrapError(err)
	}

	token := auth.Login(user)

	return loginOk(ctx, headers, &loginResponse{Token: token})
}

func loginOk(ctx _context.EchoContext, headers map[string]string, response *loginResponse) error {
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}

	return ctx.JSON(http.StatusOK, response)
}
