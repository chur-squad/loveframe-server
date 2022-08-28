package api_handler

import (
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/friends"
	api_param "github.com/chur-squad/loveframe-server/handler/api/param"
	"github.com/labstack/echo/v4"
	"net/http"
)

type userDetailResponse struct {
	Email      string `json:"email"`
	InviteCode string `json:"inviteCode"`
}

func (h *Handler) UserDetail(c echo.Context) error {
	ctx := c.(_context.EchoContext)
	param, err := api_param.GenerateUserParam(ctx)
	if err != nil {
		return userDetailError(ctx, _error.WrapError(err))
	}

	jwt, err := h.parent.Jwt.GenerateUserJwt(param.Jwt)
	if err != nil {
		return userDetailError(ctx, _error.WrapError(err))
	}

	user, err := h.parent.Mysql.UserById(ctx, jwt.ID)
	if err != nil {
		return userDetailError(ctx, _error.WrapError(err))
	}

	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlAllowMethods:  "GET, POST, OPTIONS",
		echo.HeaderAccessControlRequestMethod: "*",
		echo.HeaderContentType:                "application/json",
	}

	body := &userDetailResponse{
		Email:      user.Email,
		InviteCode: friends.Encode(user.Id).Data,
	}

	return userDetailOk(ctx, headers, body)
}

func userDetailOk(ctx _context.EchoContext, headers map[string]string, response *userDetailResponse) error {
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}

	return ctx.JSON(http.StatusOK, response)
}

func userDetailError(ctx _context.EchoContext, err error) error {
	return err
}
