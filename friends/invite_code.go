package friends

import (
	"github.com/chur-squad/loveframe-server/internal"
	"strconv"
)

type EncryptedInviteCode struct {
	Data string
}

func Encode(id int64) (code *EncryptedInviteCode) {
	return &EncryptedInviteCode{Data: strconv.FormatInt(id, 10)}
}

func Decode(code *EncryptedInviteCode) (id int64) {
	id, _ = internal.InterfaceToInt64(code.Data)
	return id
}
