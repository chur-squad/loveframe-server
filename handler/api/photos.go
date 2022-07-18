package api_handler

import (
	_error "github.com/chur-squad/loveframe-server/error"
	_context "github.com/chur-squad/loveframe-server/context"
	api_param "github.com/chur-squad/loveframe-server/handler/api/param"
	"github.com/labstack/echo/v4"
)

const (
	headerCacheControl = "Cache-Control"
)

// Manifest is serving a content manifest file which customized manually.
// test request
//http://localhost:8080/api/photos?jwt=ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SmxlSEFpT2pFMk5UZ3lNall3T0RNc0ltWnlhV1Z1WkVsa0lqb2lhRzl1WjJKcGJpSXNJbkJoZEhSbGNtNGlPaUl2YW1GbGFIbDFiaTkwWlhOMExtcHdaV2NpZlEuTEt0RndWM1hRYUwwYUE1S2xaYWRMa3hrYU5EWlZ4SkN2dWdxbUFyRFJSYw
func (h *Handler) Photos(c echo.Context) error {
	ctx := c.(_context.EchoContext)
	// parse parameters
	param, err := api_param.GenerateUserParam(ctx)
	if err != nil {
		return photoError(ctx, _error.WrapError(err))
	}
	// parse jwt struct for manifest
	jwt, err := h.parent.Jwt.GenerateUserJwt(param.Jwt)
	
	if err != nil {
		return photoError(ctx, _error.WrapError(err))
	}

	photo, err := h.parent.Photo.GetPhotoFromCdn(ctx, jwt)
	if err != nil {
		return photoError(ctx, _error.WrapError(err))
	}
	// add access control headers
	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlRequestMethod: "*",
		headerCacheControl: "max-age=90000, public",
		echo.HeaderContentType: "image/jpeg",
	}
	
	return photoOK(ctx, headers, photo)
}

func photoError(ctx _context.EchoContext, err error) error {
	// if a request returns http 4xx or 5xx, cloudfront caches about 5 minutes.
	// maybe if http 400, 403, 412, 415 status with cache-control header returns, it's able to control cache time.
	return err
}

func photoOK(ctx _context.EchoContext, headers map[string]string, body []byte) error {
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}
	if _, err := ctx.Response().Write(body); err != nil {
		return _error.WrapError(err)
	}
	return nil
}
