package orderbook

type TimeInForce rune

const (
	TimeInForceDay            TimeInForce = '0'
	TimeInForceGoodTillCancel             = '1'
	TimeInForceUnknown                    = '?'
)

func TimeInForceConv(str string) TimeInForce {
	switch str {
	case "Day":
		return TimeInForceDay
	case "GoodTillCancel":
		return TimeInForceGoodTillCancel
	}
	return TimeInForceUnknown
}
