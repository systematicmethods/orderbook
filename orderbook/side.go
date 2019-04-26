package orderbook

import "strings"

type Side int

const (
	SideSell    Side = -1
	SideBuy          = 1
	SideUnknown      = 2
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
