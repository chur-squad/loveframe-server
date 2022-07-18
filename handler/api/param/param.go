package api_param

import (
	"encoding/base64"
	"strings"
	"github.com/chur-squad/loveframe-server/internal"
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
)

type (
	UserParam struct {
		Jwt    string
		Format string
	}
)

// GenerateManifestParam is to return param for contents manifest.
func GenerateUserParam(ctx _context.EchoContext) (param UserParam, err error) {
	if ctx == nil {
		err = internal.ErrInvalidParams
		return
	}

	userParam, suberr := ctx.GetParam()
	if suberr != nil {
		err = _error.WrapError(suberr)
		return
	}

	// check parameters is empty or not.
	encodedJwt, _ := internal.InterfaceToString(userParam["jwt"])
	if encodedJwt == "" {
		err = _error.WrapError(internal.ErrInvalidParams)
		return
	}

	// [warning] it must be the unpadded base64 encoding string. (padding is `=` string)
	// RawURLEncoding deals with unpadded base64 string.
	rawJwt, suberr := base64.RawURLEncoding.DecodeString(encodedJwt)
	if suberr != nil {
		err = _error.WrapError(suberr)
		return
	}

	param.Jwt = string(rawJwt)
	seps := strings.Split(ctx.Request().URL.Path, "/")
	param.Format = strings.Split(seps[len(seps)-1], ".")[1]
	return
}
