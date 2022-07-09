package api_handler

import (
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/labstack/echo/v4"
)

const (
	headerCacheControl = "Cache-Control"
)

// Manifest is serving a content manifest file which customized manually.
func (h *Handler) Photos(c echo.Context) error {
	ctx := c

	photo, err := h.parent.Photo.GetPhotoFromCdn(ctx)
	if err != nil {
		return photoError(ctx, _error.WrapError(err))
	}
	// add access control headers
	headers := map[string]string{
		echo.HeaderAccessControlAllowOrigin:   "*",
		echo.HeaderAccessControlRequestMethod: "*",
	}
	
	return photoOK(ctx, headers, photo)
}

// https://docs.aws.amazon.com/AmazonCloudFront/latest/DeveloperGuide/HTTPStatusCodes.html#HTTPStatusCodes-cached-errors-general
func photoError(ctx echo.Context, err error) error {
	// if a request returns http 4xx or 5xx, cloudfront caches about 5 minutes.
	// maybe if http 400, 403, 412, 415 status with cache-control header returns, it's able to control cache time.
	return err
}

func photoOK(ctx echo.Context, headers map[string]string, body []byte) error {
	for k, v := range headers {
		ctx.Response().Header().Set(k, v)
	}
	if _, err := ctx.Response().Write(body); err != nil {
		return _error.WrapError(err)
	}
	return nil
}
