package main

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/chur-squad/loveframe-server/internal"
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	corsAllowOrigins = []string{
		"https://localhost:8080", "http://localhost:8080",
	}
	corsAllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
)

// addMiddleware is to set middleware chains on Echo.
func addMiddleware(e *echo.Echo) error {
	if e == nil {
		return internal.ErrInvalidParams
	}
	var err error

	// set debug
	e.Debug = true // when service deployed, need to change 
	// set logging
	e.HideBanner = true
	e.HidePort = true
	
	// this timeout will wait until read to request body
	e.Server.ReadTimeout = 20 * time.Second
	// this timeout will wait until between read a body and write a response
	e.Server.WriteTimeout = 40 * time.Second

	// get middleware chains which are required essentially.
	requiredChains, err := requiredMiddlewareChain()
	if err != nil {
		return err
	}
	e.Use(requiredChains...)

	return nil
}

// requiredMiddlewareChain is middleware chains that must be executed for everything requests.
func requiredMiddlewareChain() ([]echo.MiddlewareFunc, error) {
	var chain []echo.MiddlewareFunc

	// set top recover handler on panic
	chain = append(chain, middleware.RecoverWithConfig(
		middleware.RecoverConfig{StackSize: 1 << 10},
	))

	// remove slash
	chain = append(chain, middleware.RemoveTrailingSlash())

	// redirect from www.blah.com to blah.com
	chain = append(chain, middleware.NonWWWRedirect())

	// rewrite path
	chain = append(chain, middleware.Rewrite(map[string]string{
		"/static/*": "/public/$1",
	}))

	// set Cross-Origin Resource Sharing
	chain = append(chain, middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowCredentials: true,
			AllowOrigins:     corsAllowOrigins,
			AllowMethods:     corsAllowMethods,
		}))

	// set "X-Request-ID" header(it take purpose for analysis)
	// "X-Request-ID" is created from a client and returns response including this header.
	chain = append(chain, middleware.RequestIDWithConfig(
		middleware.RequestIDConfig{
			Generator: func() string { return "" },
		},
	))

	// set custom context
	chain = append(chain, func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// initializecontext without context canceled.
			mctx := _context.NewEchoContextWithoutCancel(c)
			mctx.WithContext("start_time", time.Now())

			// set deadline(20 seconds)
			ctx, cancel := context.WithTimeout(mctx.GetContext(), 20*time.Second)
			defer cancel()

			// set new context
			mctx.SetContext(ctx)

			return h(mctx)
		}
	})

	// set main middleware
	chain = append(chain, func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			mctx := c.(_context.EchoContext)
			startTime := mctx.ValueContext("start_time").(time.Time)

			// set ip
			mctx.SetExtra("ip", mctx.RealIP())
			//IP setting

			// call handler
			ckerr := h(mctx)

			// calculate execution time
			stopTime := time.Now()
			duration := stopTime.Sub(startTime)

			// create log

			status := http.StatusOK
			logMap := mctx.GetLog()
			logMap["latency"] = strconv.FormatInt(int64(duration), 10)
			logMap["latency_human"] = duration.String()

			if ckerr != nil {
				if code, httpErr := _error.ExtractEchoHttpError(ckerr); httpErr == nil {
					status = http.StatusInternalServerError
				} else {
					status = code
				}
				logMap["status"] = status
				logMap["message"] = ckerr.Error()
				mctx.SetExtra("error", ckerr.Error())
	
			} else {
				status = mctx.Response().Status
				logMap["status"] = status
			}

			mctx.SetExtra("status", strconv.Itoa(status))
			mctx.SetExtra("duration", strconv.FormatFloat(float64(duration.Nanoseconds())/1e9, 'f', -1, 64))
			mctx.Logger().Infoj(logMap)

			// if this request doesn't belong to API group, returning immediately.
			if !strings.HasPrefix(mctx.Request().URL.Path, "/api") {
				return ckerr
			}
			return ckerr
		}
	})

	return chain, nil
}