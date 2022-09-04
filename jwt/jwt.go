package jwt

import (
	"encoding/base64"
	"fmt"
	"github.com/chur-squad/loveframe-server/mysql"
	"reflect"
	"strconv"
	"time"

	_error "github.com/chur-squad/loveframe-server/error"
	"github.com/chur-squad/loveframe-server/internal"
	"github.com/dgrijalva/jwt-go"
	"github.com/gjbae1212/go-wraperror"
)

// ParseJwtByHMAC256 parses a token and validates.
func ParseJwtByHMAC256(data string, hmacSalt []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		//jwt.Parse method will be return token, token include header + payload + signature
		if token.Method.Alg() != hs256 {
			return nil, _error.WrapError(internal.ErrJwtAlgNotHMAC256)
		}
		return hmacSalt, nil
	})
	if err != nil {
		return nil, _error.WrapError(err)
	}
	if token == nil || !token.Valid {
		return nil, _error.WrapError(internal.ErrJwtInvalid)
	}
	return token, nil
}

// NewJwtStringByHMAC256 creates JWT string encrypted with salt.
func NewJwtStringByHMAC256(claim jwt.Claims, hmacSalt []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(hmacSalt)
}

// unmarshalJwt is unmarshaling jwt.Token to a struct like PingJwt, ManifestJwt.
// version field is required in everything JWT.
// bypass argument ignores to check jwt validation.

func unmarshalJwt(claims jwt.MapClaims, bypass bool, v interface{}) error {
	if claims == nil || v == nil {
		return _error.WrapError(internal.ErrInvalidParams)
	}

	// extract jwt version
	var currentVersion int64
	var err error
	if _, ok := claims["version"]; ok {
		currentVersion, err = internal.InterfaceToInt64(claims["version"])
		//버전 정보 체크
		if err != nil {
			return _error.WrapError(err)
		}
	}

	value := reflect.ValueOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		v := value.Field(i)
		vt := value.Type().Field(i)
		if tag, ok := vt.Tag.Lookup("name"); ok {
			//conventional format
			required := vt.Tag.Get("tag")
			elem, eok := claims[tag]

			if required == "required" && !eok && !bypass {
				minStr, mok := vt.Tag.Lookup("min")
				if !mok {
					return _error.WrapError(internal.ErrJwtInvalid)
				}
				// if a value is empty, an error is raising
				// when a tag is `required` and `current_version` is the same or higher than `min`.
				min, _ := strconv.ParseInt(minStr, 10, 64)
				if currentVersion >= min {
					chainErr := wraperror.Error(fmt.Errorf("[err] unmarshalJwt %s equal and higher version (%d != %d)",
						tag, min, currentVersion)).Wrap(internal.ErrJwtInvalid)
					return _error.WrapError(chainErr)
				}
			} else if eok {
				if elem == nil {
					continue
				}
				switch elem.(type) {
				case int:
					v.Set(reflect.ValueOf(int64(elem.(int))))
				case int32:
					v.Set(reflect.ValueOf(int64(elem.(int32))))
				case int64:
					v.Set(reflect.ValueOf(elem))
				case float32:
					v.Set(reflect.ValueOf(int64(elem.(float32))))
				case float64:
					v.Set(reflect.ValueOf(int64(elem.(float64))))
				case []interface{}:
					e, err := internal.InterfaceSliceToStringSlice(elem.([]interface{}))
					if err != nil {
						return _error.WrapError(err)
					}
					v.Set(reflect.ValueOf(e))
				default:
					v.Set(reflect.ValueOf(elem))
				}
			}
		}
	}
	return nil
}

func CreateJwt(user *mysql.User) string {

	baseTime := time.Now().Add(100 * time.Hour)

	Claims := jwt.MapClaims{
		"id":        user.Id,
		"name":      user.Name,
		"friend_id": user.FriendId,
		"exp":       baseTime.Unix(),
		"pattern":   "/photos/jaehyun/test.jpg",
	}
	token, _ := NewJwtStringByHMAC256(Claims, []byte("loveframe"))
	jwt := base64.RawURLEncoding.EncodeToString([]byte(token))

	return jwt
}
