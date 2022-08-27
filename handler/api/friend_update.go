package api_handler

import (
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
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

	request := new(addFriendRequest)
	if err = c.Bind(request); err != nil {
		return _error.WrapError(err)
	}

	err = h.parent.Friend.AddFriend(ctx, jwt, request.EncryptedInviteCode)
	if err != nil {
		return _error.WrapError(err)
	}

	return nil
}
