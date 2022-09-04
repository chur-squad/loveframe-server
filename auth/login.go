package auth

import (
	"github.com/chur-squad/loveframe-server/jwt"
	"github.com/chur-squad/loveframe-server/mysql"
)

func Login(user *mysql.User) string {
	return jwt.CreateJwt(user)
}
