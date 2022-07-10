package main

import (
	"fmt"
	"net/http"
	"github.com/chur-squad/loveframe-server/handler"
	api_handler "github.com/chur-squad/loveframe-server/handler/api"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/labstack/echo/v4"
)

// addRoute is to set route rule on Echo.
func addRoute(e *echo.Echo, h *handler.Handler) error {
	if e == nil || h == nil {
		return _error.WrapError(_error.ErrInvalidParams)
	}

	// error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// extract http error
		if status, httpErr := _error.ExtractEchoHttpError(err); httpErr == nil {
			c.JSON(http.StatusInternalServerError, map[string]string{})
		} else {
			switch httpErr.Message.(type) {
			case string:
				c.JSON(status, map[string]string{"message": fmt.Sprintf("%v", httpErr.Message)})
			default:
				c.JSON(status, httpErr.Message)
			}
		}
	}

	// static file handler
	e.Static("/", "public")

	apiH, err := api_handler.NewHandler(h)
	if err != nil {
		return _error.WrapError(err)
	}

	// api
	apiChains, err := groupApiMiddlewareChain()
	if err != nil {
		return _error.WrapError(err)
	}
	apiGroup := e.Group("/api", apiChains...)
	// making api chain by grouping

	//apiGroup.GET("/friends/*", apiH.Friends)
	// photos API need jwt from context
	// when jwt query set in context?
	apiGroup.GET("/photos/*", apiH.Photos)

	return nil
}

// groupApiMiddlewareChain is middleware chains that must execute when requests are incoming to have /api/ URL path.
func groupApiMiddlewareChain() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc
	return chain, nil
}