package api_param

import (
	"encoding/base64"
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/internal"
	"strings"
)

type (
	UserParam struct {
		Jwt    string
		Format string
	}
)

// GenerateUserParam is to return param for user info.
func GenerateUserParam(ctx _context.EchoContext) (param UserParam, err error) {
	if ctx == nil {
		err = internal.ErrInvalidParams
		return
	}

	// check parameters is empty or not.
	header := ctx.GetHeader()
	authorization, _ := internal.InterfaceToString(header["Authorization"])
	if authorization == "" {
		err = _error.WrapError(internal.ErrInvalidParams)
		return
	}

	encodedJwt := strings.Split(authorization, "Bearer ")[1]

	// [warning] it must be the unpadded base64 encoding string. (padding is `=` string)
	// RawURLEncoding deals with unpadded base64 string.
	rawJwt, suberr := base64.RawURLEncoding.DecodeString(encodedJwt)
	if suberr != nil {
		err = _error.WrapError(suberr)
		return
	}

	param.Jwt = string(rawJwt)
	return
}
