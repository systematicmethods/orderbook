package orderbook

import "strings"

type Side rune

const (
	SideBuy     Side = '1'
	SideSell         = '2'
	SideUnknown      = 'u'
)

func SideConv(side string) Side {
	switch strings.ToLower(side) {
	case "sell":
		return SideSell
	case "buy":
		return SideBuy
	}
	return SideUnknown
}

func SideToString(side Side) string {
	switch side {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	}
	return "Unknown"
}
