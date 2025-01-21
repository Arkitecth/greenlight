package data

import (
	"fmt"
	"strconv"
)

type Runtime int32


func (r Runtime) MarshalJSON() ([]byte, error) {
	js := fmt.Sprintf("%d minutes", r)

	quotedJSONValue := strconv.Quote(js)

	return []byte(quotedJSONValue), nil
}


