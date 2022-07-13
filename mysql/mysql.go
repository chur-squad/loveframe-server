package mysql

import (
	_error "github.com/chur-squad/loveframe-server/error"
	gomysql "github.com/gjbae1212/go-sql/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Mysql interface{
	UserModel
}

type connector struct {
	loveframeDB *gorm.DB
}

// NewMysql returns a connector to have been implementing interface(Mysql).
func NewMysql(dsn string, tries int) (Mysql, error) {
	conn, err := gomysql.NewConnector(dsn, tries)
	if err != nil {
		return nil, _error.WrapError(err)
	}

	// trying connect to mysql
	if err := conn.Connect(); err != nil {
		return nil, _error.WrapError(err)
	}

	// getting *sql.DB
	db, err := conn.DB()
	if err != nil {
		return nil, _error.WrapError(err)
	}

	gdb, err := gorm.Open(mysql.New(mysql.Config{DriverName: "mysql", Conn: db}), &gorm.Config{})
	if err != nil {
		return nil, _error.WrapError(err)
	}

	return &connector{loveframeDB: gdb}, nil
}

