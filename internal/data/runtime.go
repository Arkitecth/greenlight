package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)



var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int32


func (r Runtime) MarshalJSON() ([]byte, error) {
	js := fmt.Sprintf("%d minutes", r)

	quotedJSONValue := strconv.Quote(js)

	return []byte(quotedJSONValue), nil
}



func (r * Runtime) UnmarshalJSON(jsonValue []byte ) error {
	unqoutedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	splitString := strings.Split(unqoutedJSONValue, " ")

	if len(splitString) != 2 || splitString[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}


	i, err := strconv.ParseInt(splitString[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(i)

	return nil
}

