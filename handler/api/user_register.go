package api_handler

import (
	"encoding/json"
	"io"

	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/labstack/echo/v4"
)


type userInfo struct{
	ID			int64
	Name 		string
	FriendID	int64
}

func (h *Handler) UserRegister(c echo.Context) error {
	ctx := c.(_context.EchoContext)

	bytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}

	var user userInfo
	json.Unmarshal([]byte(bytes), &user)

	err = h.parent.Mysql.AddUser(user.ID, user.Name, user.FriendID)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}
	// add access control headers
	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlAllowMethods:	"GET, POST, OPTIONS",
		echo.HeaderAccessControlRequestMethod: "*",
		echo.HeaderContentType: "application/json",
	}

	return userRegisterOK(ctx, headers)	
}

func userError(ctx _context.EchoContext, err error) error {
	return err
}

func userRegisterOK(ctx _context.EchoContext, headers map[string]string) error{
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}
	return nil
}
