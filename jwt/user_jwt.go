package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/chur-squad/loveframe-server/internal"
	_error "github.com/chur-squad/loveframe-server/error"
)

// UsertJwt is user struct for server appropriate photo.
type UserJwt struct {
	//Created			string `name:"created", tag:"required", min"0"` 
	//LastUpdated		string `name:"last_updated, tag:"required", min"0"`
	ID				int64  	`name:"id", tag:"required", min"0"`
	Name			string  `name:"name", tag:"required", min"0"`
	FriendID		int64	`name:"friend_id" tag:"required" min:"0"`
	Exp				int64	`name:"exp" tag:"required" min:"0"`
	Pattern			string	`name:"pattern" tag:"required" min:"0"`
}

func (m *manager) newUserJwt(jwtToken *jwt.Token) (UserJwt, error) {
	// func (class) func_name (param) (return)
	if jwtToken == nil {
		return UserJwt{}, _error.WrapError(internal.ErrInvalidParams)
	}

	claims := jwtToken.Claims.(jwt.MapClaims)
	userJwt := &UserJwt{}
	if err := unmarshalJwt(claims, false, userJwt); err != nil {
		return UserJwt{}, _error.WrapError(err)
	}

	return *userJwt, nil
}
