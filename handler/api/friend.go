package api_handler

import (
	"encoding/json"
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/friends"
	"github.com/labstack/echo/v4"
	"io"
)

type addFriendRequest struct {
	EncryptedInviteCode string `json:"inviteCode"`
}

func (h *Handler) AddFriend(c echo.Context) error {
	ctx := c.(_context.EchoContext)

	bytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return userError(ctx, _error.WrapError(err))
	}

	var me userInfo
	json.Unmarshal(bytes, &me)

	request := new(addFriendRequest)
	if err = c.Bind(request); err != nil {
		return _error.WrapError(err)
	}

	code := &friends.EncryptedInviteCode{Data: request.EncryptedInviteCode}
	err = h.parent.Friend.AddFriend(ctx, me.ID, code)
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
