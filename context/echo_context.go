package context

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	EchoContext interface {
		echo.Context

		GetContext() context.Context
		SetContext(ctx context.Context)
		WithContext(key, val interface{})
		ValueContext(key interface{}) interface{}
		GetLog() map[string]interface{}
		GetHeader() EchoHeader
		GetParam() (EchoParam, error)
		GetExtra() EchoInfo
		SetExtra(k, s string)
	}

	EchoParam  map[string]interface{}
	EchoHeader map[string]string
	EchoInfo   map[string]string

	echoContext struct {
		echo.Context

		param  EchoParam
		header EchoHeader
		extra  EchoInfo
	}
)

// GetContext returns a context on a request object.
func (c *echoContext) GetContext() context.Context {
	return c.Request().Context()
}

// SetContext set a new context into a request object.
func (c *echoContext) SetContext(ctx context.Context) {
	c.SetRequest(c.Request().WithContext(ctx))
}

// WithContext makes a context using args(key, value), and then it's to set into a request object.
func (c *echoContext) WithContext(key, val interface{}) {
	ctx := c.GetContext()
	c.SetContext(context.WithValue(ctx, key, val))
}

// ValueContext returns a value matched a key on a request context
func (c *echoContext) ValueContext(key interface{}) interface{} {
	return c.GetContext().Value(key)
}

// GetLog returns log struct.
func (c *echoContext) GetLog() map[string]interface{} {
	return map[string]interface{}{
		"host":       c.Request().Host,
		"ip":         c.RealIP(),
		"uri":        c.Request().RequestURI,
		"method":     c.Request().Method,
		"referer":    c.Request().Referer(),
		"user-agent": c.Request().UserAgent(),
	}
}

// GetHeader returns header information that is constructed a type(map[string]string).
func (c *echoContext) GetHeader() EchoHeader {
	// if param is created already.
	if c.header != nil {
		return c.header
	}

	header := EchoHeader{}
	for k, v := range c.Request().Header {
		combined := strings.Join(v, ",")
		header[k] = combined
	}
	header["Host"] = c.Request().Host
	c.header = header
	return c.header
}

// GetParam returns body param.
func (c *echoContext) GetParam() (EchoParam, error) {
	// if param is created already.
	if c.param != nil {
		return c.param, nil
	}

	// create param.
	param := map[string]interface{}{}
	if err := c.Bind(&param); err != nil {
		return nil, err
	}

	c.param = param
	return c.param, nil
}

// GetExtra returns extra information.
func (c *echoContext) GetExtra() EchoInfo {
	if c.extra != nil {
		return c.extra
	}

	info := EchoInfo{}
	c.extra = info
	return c.extra
}

// SetExtra returns extra information.
func (c *echoContext) SetExtra(k, s string) {
	extraInfo := c.GetExtra()
	extraInfo[k] = s
}

// NewEchoContextWithoutCancel returns EchoContext without context canceled.
func NewEchoContextWithoutCancel(ctx echo.Context) EchoContext {
	echoCtx := &echoContext{Context: ctx}
	echoCtx.SetContext(&WithoutCancelCtx{Context: echoCtx.GetContext()})
	return echoCtx
}

// NewEchoContext returns EchoContext.
func NewEchoContext(ctx echo.Context) EchoContext {
	echoCtx := &echoContext{Context: ctx}
	return echoCtx
}
