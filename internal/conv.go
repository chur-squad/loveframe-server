package internal

import (
	"strconv"
	_error "github.com/chur-squad/loveframe-server/error"
)

// InterfaceToInt64 converts value having interface type to value having int64.
func InterfaceToInt64(i interface{}) (int64, error) {
	if i == nil {
		return 0, _error.ErrInvalidParams
	}

	switch i.(type) {
	case int:
		return int64(i.(int)), nil
	case int64:
		return i.(int64), nil
	case int32:
		return int64(i.(int32)), nil
	case float32:
		return int64(i.(float32)), nil
	case float64:
		return int64(i.(float64)), nil
	case string:
		return strconv.ParseInt(i.(string), 10, 64)
	default:
		return 0, _error.ErrUnknown
	}
}

// InterfaceToBool converts value having interface type to value having bool.
func InterfaceToBool(i interface{}) (bool, error) {
	if i == nil {
		return false, _error.ErrInvalidParams
	}

	switch i.(type) {
	case string:
		return strconv.ParseBool(i.(string))
	case bool:
		return i.(bool), nil
	default:
		return false, _error.ErrUnknown
	}
}

// InterfaceToString converts a value having interface type to string
func InterfaceToString(i interface{}) (string, error) {
	if i == nil {
		return "", _error.ErrInvalidParams
	}

	switch i.(type) {
	case int:
		return strconv.FormatInt(int64(i.(int)), 10), nil
	case int64:
		return strconv.FormatInt(i.(int64), 10), nil
	case int32:
		return strconv.FormatInt(int64(i.(int32)), 10), nil
	case float32:
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 64), nil
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(i.(bool)), nil
	case string:
		return i.(string), nil
	default:
		return "", _error.ErrUnknown
	}
}

// InterfaceToFloat64 converts value having interface type to value having float64.
func InterfaceToFloat64(i interface{}) (float64, error) {
	if i == nil {
		return 0, _error.WrapError(_error.ErrInvalidParams)
	}

	switch i.(type) {
	case int:
		return float64(i.(int)), nil
	case int64:
		return float64(i.(int64)), nil
	case int32:
		return float64(i.(int32)), nil
	case float32:
		return float64(i.(float32)), nil
	case float64:
		return i.(float64), nil
	case string:
		return strconv.ParseFloat(i.(string), 64)
	default:
		return 0, _error.WrapError(_error.ErrInvalidParams)
	}
}


// Converts a slice of interfaces into a slice of strings.
func InterfaceSliceToStringSlice(is []interface{}) ([]string, error) {
	if is == nil {
		return nil, _error.ErrInvalidParams
	}

	ss := make([]string, 0, len(is))
	for _, i := range is {
		s, err := InterfaceToString(i)
		if err != nil {
			return nil, err
		}
		ss = append(ss, s)
	}
	return ss, nil
}
