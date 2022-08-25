package handler

import (
	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/friends"
	jwt "github.com/chur-squad/loveframe-server/jwt"
	"github.com/chur-squad/loveframe-server/mysql"
	photos "github.com/chur-squad/loveframe-server/photos"
)

type (
	// Handler is handler struct for using in Echo.
	Handler struct {
		Cfg    *Config
		Photo  photos.Manager
		Mysql  mysql.Mysql
		Jwt    jwt.Manager
		Friend friends.Manager
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

	photo, err := photos.NewManager(
		photos.WithCdnEndpoint(h.Cfg.CdnEndpoint),
	)

	mysql, err := mysql.NewMysql(h.Cfg.MysqlDSN, 2)

	friend, err := friends.NewManager(mysql)

	if err != nil {
		return nil, _error.WrapError(err)
	}
	h.Jwt = jwt
	h.Photo = photo
	h.Mysql = mysql
	h.Friend = friend

	return h, nil
}
