package mysql

import (
	"errors"
	_error "github.com/chur-squad/loveframe-server/error"
	_context "github.com/chur-squad/loveframe-server/context"
	"github.com/chur-squad/loveframe-server/internal"
	"gorm.io/gorm"
	"time"
)

type UserModel interface {
	UserById(ctx _context.EchoContext, id int64) (*User, error)
	AddUser(id int, name string) error
}

type User struct {
	Id                  int       `gorm:"primary_key;column:id"`
	Name                string    `gorm:"type:varchar(255);column:name"`
	Password            string    `gorm:"type:varchar(255);column:password"`
	PasswordSalt        string    `gorm:"type:varchar(255);column:password_salt"`
	UploadImageDomain   string    `gorm:"type:varchar(255);column:upload_image_domain"`
	FriendId        	string    `gorm:"type:varchar(255);column:friend_id"`
	CreatedAt           time.Time `gorm:"type:datetime;column:created_at"`
	UpdatedAt           time.Time `gorm:"type:datetime;column:updated_at"`
}

// TableName returns table-name, it is using by gorm when extracting table name.
func (u *User) TableName() string {
	return "user"
}

// IsEmpty returns whether empty or not.
func (u *User) IsEmpty() bool {
	return u.Id == 0
}

// addUser
func (c *connector) AddUser (Id int, Name string) error {
	//get context from request
	db := c.loveframeDB

	db.AutoMigrate(&User{})
	db.Create(&User{Id: Id, Name: Name, CreatedAt: time.Now(), UpdatedAt: time.Now()})
	
	var	checkUser User
	db.First(&checkUser, 1)

	return nil
} 

// UserById returns user by id.
func (c *connector) UserById(ctx _context.EchoContext, id int64) (*User, error) {
	if ctx == nil || id <= 0 {
		return nil, _error.WrapError(internal.ErrInvalidParams)
	}
	//get context from request
	db := c.loveframeDB.WithContext(ctx.Request().Context())

	user := &User{}

	// wrapping logic to function
	run := func() error {
		result := db.Where("id = ?", id).First(user)

		// if result don't exist.
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		// if an unknown error is raised.
		if result.Error != nil {
			return _error.WrapError(result.Error)
		}

		return nil
	}

	if err := run(); err != nil {
		return nil, _error.WrapError(err)
	}

	if user.IsEmpty() {
		return nil, nil
	} else {
		return user, nil
	}
}