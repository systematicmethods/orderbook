package orderbook

type OrderType rune

const (
	Market            OrderType = '1'
	Limit                       = '2'
	Stop                        = '3'
	OrderType_unknown           = 'U'
)

func OrderTypeConv(ordertype string) OrderType {
	switch ordertype {
	case "Limit":
		return Limit
	case "Market":
		return Market
	case "Stop":
		return Stop
	}
	return OrderType_unknown
}
