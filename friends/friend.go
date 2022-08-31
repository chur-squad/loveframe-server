package friends

import (
	_context "github.com/chur-squad/loveframe-server/context"
	"github.com/chur-squad/loveframe-server/internal"
	"github.com/chur-squad/loveframe-server/mysql"
)

type Manager interface {
	AddFriend(ctx _context.EchoContext, meId int64, inviteCode *EncryptedInviteCode) error
}

type FriendMaker struct {
	Mysql mysql.Mysql
}

func (maker *FriendMaker) AddFriend(ctx _context.EchoContext, meId int64, inviteCode *EncryptedInviteCode) error {
	friendId := Decode(inviteCode)

	m := maker.Mysql

	me, _ := m.UserById(ctx, meId)
	friend, _ := m.UserById(ctx, friendId)

	m.Connect(ctx, me, friend)

	return nil
}

func NewManager(m mysql.Mysql) (Manager, error) {
	maker := &FriendMaker{Mysql: m}

	return maker, nil
}

func decodeUserId(code string) int64 {
	id, _ := internal.InterfaceToInt64(code)
	return id
}

func encodeBy(user *mysql.User) string {
	encrypted, _ := internal.InterfaceToString(user.Id)
	return encrypted
}
