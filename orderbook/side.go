package orderbook

import "strings"

type Side int

const (
	Sell         Side = -1
	Buy               = 1
	side_unknown      = 2
)

func SideConv(side string) Side {
	switch strings.ToLower(side) {
	case "sell":
		return Sell
	case "buy":
		return Buy
	}
	return side_unknown
}
