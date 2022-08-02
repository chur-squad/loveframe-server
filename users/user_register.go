package users

import (
	"github.com/chur-squad/loveframe-server/internal"	
	_jwt "github.com/chur-squad/loveframe-server/jwt"
	_error "github.com/chur-squad/loveframe-server/error"
	_context "github.com/chur-squad/loveframe-server/context"
)

type Manager interface {
	GetUserInfo(ctx _context.EchoContext, jwt _jwt.UserJwt) error
}

type userMaker struct {
	userJwt			_jwt.UserJwt	
}

//make user struct for register and store DB

//Get UserInfo from UserJwt When user first registered
func (maker *userMaker) GetUserInfo(ctx _context.EchoContext, jwt _jwt.UserJwt) error {
	//Get User data from jwt claim
	//check current user data matches user struct format
	//need to define return value : Is it need to return User Struct?
	
	return nil
}

//Add User data to database 
/* move to api for proper dependency
func (maker *userMaker) AddUserToDB(jwt _jwt.UserJwt, mysql _handler.mysql) error {
	//Save user data to Database from jwt 
	err := mysql.AddUser(jwt.ID, jwt.Name)
	if err != nil {
		return _error.WrapError(internal.ErrDatabaseUpdate) 
	}
}
*/
func (maker *userMaker) Valid() (ok bool) {
	if (maker.userJwt == _jwt.UserJwt{}) {
		return
	}
	ok = true
	return
}
// NewManager returns a user object that is implemented Manager interface.
func NewManager() (Manager, error) {
	maker := &userMaker{}

	if !maker.Valid() {
		return nil, _error.WrapError(internal.ErrInvalidParams)
	}
	return maker, nil
}
