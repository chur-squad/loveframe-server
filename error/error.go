package error

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	wraperror "github.com/gjbae1212/go-wraperror"
	"github.com/labstack/echo/v4"
)

var (
	// ErrInvalidParams is an error type to use when passed arguments are invalid.
	ErrInvalidParams = errors.New("[err] invalid params")
	// ErrInvalidRequest is an error type to use when a request are invalid.
	ErrInvalidRequest = errors.New("[err] invalid request")
	// ErrUnknown is an error type to use when error reason doesn't know.
	ErrUnknown = errors.New("[err] unknown")
	// ErrParse is an error type to use when error reason doesn't parse.
	ErrParse = errors.New("[err] parse")
)

const maxStackTraceCount = 32

type ErrorWithStackTrace struct {
	error
	stackTrace []uintptr
}

// Sentry finds [StackTrace() []unitptr] to show stack trace.
// no more implementation is required to show stack trace in sentry web.
func (e *ErrorWithStackTrace) StackTrace() []uintptr {
	return e.stackTrace
}

// implement [Is(error) bool]
func (e *ErrorWithStackTrace) Is(target error) bool {
	return errors.Is(e.error, target)
}

// implement [As(interface{}) bool]
func (e *ErrorWithStackTrace) As(target interface{}) bool {
	return errors.As(e.error, target)
}

// implement [Unwrap() error]
func (e *ErrorWithStackTrace) Unwrap() error {
	return e.error
}

// ExtractEchoHttpError extracts *echo.HttpError object.
func ExtractEchoHttpError(err error) (status int, httpErr *echo.HTTPError) {
	chainErr := wraperror.Error(err)
	for _, e := range chainErr.Flatten() {
		if suberr, ok := e.(*echo.HTTPError); ok {
			status = suberr.Code
			httpErr = suberr
			return
		}
	}
	return
}

// WrapError wraps error.
func WrapError(err error) error {
	if err != nil {
		// Get stack trace. if err already has stack trace, restore it.
		var stackTrace []uintptr
		if stackErr, ok := err.(*ErrorWithStackTrace); ok {
			// The error is ErrorWithStackTrace. Just extract stack trace from it.
			stackTrace = stackErr.StackTrace()
			// maintain only one ErrorWithStackTrace in error chain. Unwrap ErrorWithStackTrace.
			err = stackErr.error
		} else {
			// The error is not ErrorWithStackTrace. Create new one.
			pcs := make([]uintptr, maxStackTraceCount)
			n := runtime.Callers(2, pcs)
			stackTrace = pcs[:n]
		}

		// Get program counter and line number
		pc, _, line, _ := runtime.Caller(1)
		// Get function name from program counter
		fn := runtime.FuncForPC(pc).Name()
		// Refine function name
		details := strings.Split(fn, "/")
		fn = details[len(details)-1]
		// Build chain
		chainErr := wraperror.Error(err)
		wrappedErr := chainErr.Wrap(fmt.Errorf("[err][%s:%d]", fn, line))

		// return error with stack trace
		return &ErrorWithStackTrace{error: wrappedErr, stackTrace: stackTrace}
	}
	return nil
}
