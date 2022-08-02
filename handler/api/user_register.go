package api_handler

import (
	_error "github.com/chur-squad/loveframe-server/error"
	_context "github.com/chur-squad/loveframe-server/context"
	api_param "github.com/chur-squad/loveframe-server/handler/api/param"
	"github.com/labstack/echo/v4"
)

func (h *Handler) UserRegister(c echo.Context) error {
	ctx := c.(_context.EchoContext)

	param, err := api_param.GenerateUserParam(ctx)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}

	jwt, err := h.parent.Jwt.GenerateUserJwt(param.Jwt)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}
	/*
	err = h.parent.User.GetUserInfo(ctx, jwt)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}
	*/
	err = h.parent.Mysql.AddUser(jwt.ID, jwt.Name, jwt.FriendID)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}
	// add access control headers
	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
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
