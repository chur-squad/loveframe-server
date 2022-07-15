package handler

import (
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/mysql"
	photos "github.com/chur-squad/loveframe-server/photos"
	jwt "github.com/chur-squad/loveframe-server/jwt"
)

type (
	// Handler is handler struct for using in Echo.
	Handler struct {
		Cfg			*Config
		Photo		photos.Manager
		Mysql       mysql.Mysql		
		Jwt			jwt.Manager
	}
)

// NewHandler is to create a handler object.
func NewHandler(opts ...Option) (*Handler, error) {
	h := &Handler{}
	//handelr is struct 

	mergeOpts := []Option{}
	mergeOpts = append(mergeOpts, opts...)
	for _, opt := range mergeOpts {
		opt.apply(h)
	}
	
	// make jwt manager.
	jwt, err := jwt.NewManager(
		jwt.WithUserJwtSalt([]byte(h.Cfg.UserJwtSalt)),
		jwt.WithUserSalt(h.Cfg.UserSalt),
		jwt.WithGroupSalt(h.Cfg.GroupSalt),
	)

	if err != nil {
		return nil, _error.WrapError(err)
	}

	mysql, err := mysql.NewMysql(h.Cfg.MysqlDSN, 2)

	if err != nil {
		return nil, _error.WrapError(err)
	}
	h.Mysql = mysql
	h.Jwt = jwt

	return h, nil
}
