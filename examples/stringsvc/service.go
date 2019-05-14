package stringsvc

import (
	"io"

	"github.com/shopspring/decimal"
)

type ListOfStrings []string

type StringService interface {
	Concat(a string, b *string) (*string, *string)
	Plus(a decimal.Decimal, b *decimal.Decimal) decimal.Decimal
	Prefix(x ListOfStrings) *ListOfStrings
	Yay() io.Reader
}
