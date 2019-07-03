package orderbook

type TimeInForce rune

const (
	TimeInForceDay            TimeInForce = '0'
	TimeInForceGoodTillCancel TimeInForce = '1'
	TimeInForceUnknown        TimeInForce = '?'
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
