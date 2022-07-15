package jwt

import (
	"github.com/chur-squad/loveframe-server/internal"
	_error "github.com/chur-squad/loveframe-server/error"
	"fmt"
)

const (
	hs256 = "HS256"
)


// Manager is a token manager that do dedicates validate, parse, and so on from token.
// token manager that do dedicates validate, parse, and so on from token

type Manager interface {
	// generate jwt
	GenerateUserJwt(encrypted string) (UserJwt, error)
}

type manager struct {
	userJwtSalt		[]byte
	userSalt        string
	groupSalt       string
}

// ManagerOption is an interface for Manager, it's used for dependency injection.
type ManagerOption interface {
	apply(m *manager)
}

// ManagerOptionFunc  implements a struct from ManagerOption interface.
type ManagerOptionFunc func(m *manager)

func (opt ManagerOptionFunc) apply(m *manager) { opt(m) }


// WithUserJwtSalt returns a function for setting salt for user JWT.
func WithUserJwtSalt(salt []byte) ManagerOptionFunc {
	return func(m *manager) { m.userJwtSalt = salt }
}

// WithUserSalt returns a function for setting salt for user.
func WithUserSalt(salt string) ManagerOptionFunc {
	return func(m *manager) { m.userSalt = salt }
}

// WithGroupSalt returns a function for setting salt for group.
func WithGroupSalt(salt string) ManagerOptionFunc {
	return func(m *manager) { m.groupSalt = salt }
}

// NewManager creates Manager interface.
func NewManager(opts ...ManagerOption) (Manager, error) {
	m := &manager{}

	// set default options
	mergeOpts := []ManagerOption{}
	// merge default options and arguments options
	mergeOpts = append(mergeOpts, opts...)
	// apply options
	for _, opt := range mergeOpts {
		opt.apply(m)
	}
	fmt.Print(string(m.userJwtSalt), len(m.userJwtSalt))
	if len(m.userJwtSalt) == 0  {
		return nil, _error.WrapError(internal.ErrInvalidParams)
	}

	return m, nil
}

// GenerateManifestJwt creates ManifestJwt struct which includes information for information for contents manifest.
func (m *manager) GenerateUserJwt(encrypted string) (UserJwt, error) {
	if encrypted == "" {
		return UserJwt{}, _error.WrapError(internal.ErrInvalidParams)
	}

	jwtToken, err := ParseJwtByHMAC256(encrypted, m.userJwtSalt)
	if err != nil {
		return UserJwt{}, _error.WrapError(err)
	}
	return m.newUserJwt(jwtToken)
}



