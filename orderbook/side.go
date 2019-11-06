package orderbook

import "strings"

type Side rune

const (
	SideBuy     Side = '1'
	SideSell    Side = '2'
	SideUnknown Side = 'u'
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

func (it Side) String() string {
	switch it {
	case SideBuy:
		return "buy"
	case SideSell:
		return "sell"
	}
	return "Unknown"
}
