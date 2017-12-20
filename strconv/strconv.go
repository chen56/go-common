package strconv

import (
	"fmt"
	"strconv"
	"errors"
)

func FormatIntOrString(x interface{}) (value string, err error) {
	switch t := x.(type) {
	case string:
		return t, nil
	case float64:
		return strconv.Itoa(int(t)), nil
		//... etc
	default:
		return "", errors.New(fmt.Sprintf("expected type: int | string,'%v' actual type: '%T'", x, x))
	}
}