package orderbook

import "github.com/shopspring/decimal"

func roundup(num decimal.Decimal, places int32) decimal.Decimal {
	// math.Floor(x*100)/100
	factor := decimal.NewFromFloat32(10).Pow(decimal.New(int64(places), 0))
	return num.Mul(factor).Ceil().Div(factor)
}

func rounddown(num decimal.Decimal, places int32) decimal.Decimal {
	// math.Floor(x*100)/100
	factor := decimal.NewFromFloat32(10).Pow(decimal.New(int64(places), 0))
	return num.Mul(factor).Floor().Div(factor)
}
