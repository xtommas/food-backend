package data

import (
	"fmt"
	"strconv"
)

type Price float64

func (p Price) MarshalJSON() ([]byte, error) {
	JSONValue := fmt.Sprintf("$%g", p)

	quotedJSONValue := strconv.Quote(JSONValue)

	return []byte(quotedJSONValue), nil
}
