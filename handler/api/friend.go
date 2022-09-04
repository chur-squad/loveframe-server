package api_handler

import (
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/friends"
	api_param "github.com/chur-squad/loveframe-server/handler/api/param"
	"github.com/labstack/echo/v4"
)

type addFriendRequest struct {
	EncryptedInviteCode string `json:"inviteCode"`
}

func (h *Handler) AddFriend(c echo.Context) error {
	ctx := c.(_context.EchoContext)

	param, err := api_param.GenerateUserParam(ctx)
	if err != nil {
		return _error.WrapError(err)
	}

	jwt, err := h.parent.Jwt.GenerateUserJwt(param.Jwt)
	if err != nil {
		return _error.WrapError(err)
	}

	request := new(addFriendRequest)
	if err = c.Bind(request); err != nil {
		return _error.WrapError(err)
	}

	code := &friends.EncryptedInviteCode{Data: request.EncryptedInviteCode}
	err = h.parent.Friend.AddFriend(ctx, jwt.ID, code)
	if err != nil {
		return _error.WrapError(err)
	}

	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlAllowMethods:  "GET, POST, OPTIONS",
		echo.HeaderAccessControlRequestMethod: "*",
		echo.HeaderContentType:                "application/json",
	}

	return friendConnectOk(ctx, headers)
}

func friendConnectOk(ctx _context.EchoContext, headers map[string]string) error {
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}
	return nil
}
