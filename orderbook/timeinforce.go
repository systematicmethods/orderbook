package orderbook

type TimeInForce rune

const (
	TimeInForceDay               TimeInForce = '0'
	TimeInForceGoodTillCancel    TimeInForce = '1'
	TimeInForceImmediateOrCancel TimeInForce = '3'
	TimeInForceFillOrKill        TimeInForce = '4'
	TimeInForceGoodTillDate      TimeInForce = '6'
	TimeInForceGoodForTime       TimeInForce = 'A'
	TimeInForceGoodForAuction    TimeInForce = 'B'
	TimeInForceUnknown           TimeInForce = '?'
)

func TimeInForceConv(str string) TimeInForce {
	switch str {
	case "Day":
		return TimeInForceDay
	case "GoodForAuction":
		return TimeInForceGoodForAuction
	case "GoodTillCancel":
		return TimeInForceGoodTillCancel
	}
	return TimeInForceUnknown
}
