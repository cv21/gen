//+build ignore
//go:generate gen

package stringsvc

import "github.com/shopspring/decimal"

type StringService interface {
	Concat(a string, b *string) (*string, *string)
	Plus(a decimal.Decimal, b *decimal.Decimal) decimal.Decimal
}
