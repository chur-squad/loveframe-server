package api_handler

import (
	"encoding/json"
	"io"

	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/labstack/echo/v4"
)

type photoInfo struct {
	Photo		byte
	User       	userInfo // Client needs to return this struct
	Pattern 	string // it will be produce in serverside by userInfo
}

func (h *Handler) PhotoUpdate(c echo.Context) error {
	ctx := c.(_context.EchoContext)

	bytes, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return photoUpdateError(ctx, _error.WrapError(err))
	}

	var photo photoInfo
	json.Unmarshal([]byte(bytes), &photo)

	//photoInfo makes S3 pattern 
	//if not User Directory exist, make User Directory at S3
	//Add current photo, decide policy about former photo
	//select S3 connecting package

	// add access control headers
	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlAllowMethods:  "GET, POST, OPTIONS",
		echo.HeaderAccessControlRequestMethod: "*",
		echo.HeaderContentType:                "application/json",
	}

	return photoUpdateOK(ctx, headers)
}

func photoUpdateError(ctx _context.EchoContext, err error) error {
	// if a request returns http 4xx or 5xx, cloudfront caches about 5 minutes.
	// maybe if http 400, 403, 412, 415 status with cache-control header returns, it's able to control cache time.
	return err
}

func photoUpdateOK(ctx _context.EchoContext, headers map[string]string) error {
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}
	return nil
}
