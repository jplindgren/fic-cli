package model

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Stock struct {
	Ticker string
	Price  decimal.Decimal
	Target decimal.Decimal
}

func (s Stock) ToString() string {
	return fmt.Sprintf("%s, price: %s, target %s", s.Ticker, s.Price.StringFixed(2), s.Target.StringFixed(2))
}

func (s Stock) Ratio() float64 {
	fPrice, _ := s.Price.Float64()
	fTarget, _ := s.Target.Float64()
	ratio := fPrice / fTarget
	return ratio
}

func (s Stock) IsRecommended() (bool, string) {
	ratio := s.Ratio()

	if ratio < 0.5 {
		return true, Excellent
	} else if ratio < 0.9 {
		return true, Good
	} else if ratio < 1.1 {
		return false, InTarget
	} else {
		return false, NotRecommended
	}
}

var Excellent = "much lower than target"
var Good = "below the target"
var InTarget = "at target"
var NotRecommended = "higher than target"
