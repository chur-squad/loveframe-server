package internal

import (
	"errors"
)

var (
	ErrUnauthenticated = errors.New("[err] unauthenticated")
	// ErrInvalidParams is an error type to use when passed arguments are invalid.
	ErrInvalidParams = errors.New("[err] invalid params")
	// ErrInvalidJwtPayload is an error type to use when passed ping payload is invalid.
	ErrInvalidJwtPayload = errors.New("[err] invalid jwt payload")
	// ErrUnknown is an error type to use when error reason doesn't know.
	ErrUnknown = errors.New("[err] unknown")
	// ErrParse is an error type to use when error reason doesn't parse.
	ErrParse = errors.New("[err] parse")
	// ErrConflict occurs when response must be redirected with status of 409.
	ErrConflict = errors.New("[err] conflict status 409")
	// ErrJwtInvalid occurs when this token is invalid.
	ErrJwtInvalid = errors.New("[err] this token is invalid")
	// ErrJwtAlgNotHMAC256 occurs when this token doesn't create to using hs256.
	ErrJwtAlgNotHMAC256 = errors.New("[err] this token doesn't create to using hmac256")
	ErrUserNotFound = errors.New("[err] user not found")
	// ErrReqInvalidArguments occurs when request arguments is invalid or not exist.
	ErrReqInvalidArguments            = errors.New("[err] invalid arguments for requests")
	ErrUnsupportedS3BucketKey = errors.New("[err] unsupported s3 bucket key")
	ErrDatabaseUpdate = errors.New("[err] database prcoess is failed")
)
