package photos

import (
	_context "github.com/chur-squad/loveframe-server/context"
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/internal"
	_jwt "github.com/chur-squad/loveframe-server/jwt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

type Extension string

const (
	JPG  Extension = "jpg"
	JPEG Extension = "jpeg"
	HEIC Extension = "heic"
)

const (
	contentTypeImage = "image"
)

var (
	defaultOpts = []Option{
		WithCdnClient(&http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 1 * time.Minute,
				}).DialContext,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
				ResponseHeaderTimeout: 9 * time.Second,
				IdleConnTimeout:       5 * time.Minute,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   10,
			},
			Timeout: 10 * time.Second,
		}),
	}
)

type Manager interface {
	GetPhotoFromCdn(ctx _context.EchoContext, jwt _jwt.UserJwt) ([]byte, error)
}

type photoMaker struct {
	// cdn
	cdnClient   *http.Client
	cdnEndpoint string
}

// Valid checks this object is valid or not.
func (maker *photoMaker) Valid() (ok bool) {
	if maker.cdnClient == nil || maker.cdnEndpoint == "" {
		return
	}
	ok = true
	return
}

// GetPhotoFromCdn @photoMaker reads original photo from cdn
func (maker *photoMaker) GetPhotoFromCdn(ctx _context.EchoContext, jwt _jwt.UserJwt) ([]byte, error) {
	// get cdn endpoint
	endpoint, err := maker.getCdnEndpoint()
	if err != nil {
		return nil, _error.WrapError(err)
	}

	// get key
	key, err := maker.getKeyFromCtxAndJwt(ctx, jwt)
	if err != nil {
		return nil, _error.WrapError(err)
	}

	//get photo
	req, err := http.NewRequest(http.MethodGet, endpoint+"/"+key, nil) // method Get
	if err != nil {
		return nil, _error.WrapError(err)
	}

	resp, err := maker.requestCdn(ctx, req)
	if err != nil {
		return nil, _error.WrapError(err)
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		return nil, _error.WrapError(_error.ErrUnknown)
	}

	// read image to byte
	bys, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, _error.WrapError(err)
	}

	return bys, nil
}

// getCdnEndpointFromJwt @photoMaker returns cdn endpoint for given jwt
func (maker *photoMaker) getCdnEndpoint() (string, error) {
	return maker.cdnEndpoint, nil
}

// getKeyFromCtxAndJwt @photoMaker returns key for given context and jwt
func (maker *photoMaker) getKeyFromCtxAndJwt(ctx _context.EchoContext, jwt _jwt.UserJwt) (string, error) {
	// check if request path and jwt path match
	var pathPattern, pathExtension, jwtPattern, jwtExtension string
	if pathSeps := strings.Split(strings.TrimPrefix(ctx.Request().URL.Path, "/api/photos/"), "/"); len(pathSeps) >= 1 {
		pathPattern = strings.Join(pathSeps[:len(pathSeps)-1], "/")
		if seps := strings.Split(pathSeps[len(pathSeps)-1], "."); len(seps) == 2 {
			pathExtension = seps[len(seps)-1]
		}
	}
	if pathSeps := strings.Split(strings.TrimPrefix(jwt.Pattern, "/"), "/"); len(pathSeps) >= 1 {
		jwtPattern = strings.Join(pathSeps[:len(pathSeps)-1], "/")
		if seps := strings.Split(pathSeps[len(pathSeps)-1], "."); len(seps) == 2 {
			jwtExtension = seps[len(seps)-1]
		}
	}
	if pathPattern != jwtPattern || pathExtension != jwtExtension {
		return "", _error.WrapError(internal.ErrUnsupportedS3BucketKey)
	}
	return strings.TrimPrefix(jwt.Pattern, "/"), nil
}

func (maker *photoMaker) requestCdn(ctx _context.EchoContext, req *http.Request) (*http.Response, error) {
	resp, err := maker.cdnClient.Do(req)
	return resp, err
}

// NewManager returns a photo object that is implemented Manager interface.
func NewManager(opts ...Option) (Manager, error) {
	maker := &photoMaker{}

	mergeOpts := []Option{}
	mergeOpts = append(mergeOpts, defaultOpts...)
	mergeOpts = append(mergeOpts, opts...)

	for _, opt := range mergeOpts {
		opt.apply(maker)
	}
	if !maker.Valid() {
		return nil, _error.WrapError(internal.ErrInvalidParams)
	}
	return maker, nil
}
