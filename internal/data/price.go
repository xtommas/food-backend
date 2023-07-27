package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidPriceFormat = errors.New("invalid price format")

type Price float64

func (p Price) MarshalJSON() ([]byte, error) {
	JSONValue := fmt.Sprintf("$%g", p)

	quotedJSONValue := strconv.Quote(JSONValue)

	return []byte(quotedJSONValue), nil
}

func (p *Price) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidPriceFormat
	}

	number := strings.TrimPrefix(unquotedJSONValue, "$")

	i, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return ErrInvalidPriceFormat
	}

	*p = Price(i)

	return nil
}
