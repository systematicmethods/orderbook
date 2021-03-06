package orderbookex

type OrderType rune

const (
	OrderTypeMarket  OrderType = '1'
	OrderTypeLimit   OrderType = '2'
	OrderTypeStop    OrderType = '3'
	OrderTypeUnknown OrderType = 'U'
)

func OrderTypeConv(ordertype string) OrderType {
	switch ordertype {
	case "Limit":
		return OrderTypeLimit
	case "Market":
		return OrderTypeMarket
	case "Stop":
		return OrderTypeStop
	}
	return OrderTypeUnknown
}

func (it OrderType) String() string {
	switch it {
	case OrderTypeLimit:
		return "Limit"
	case OrderTypeMarket:
		return "Market"
	case OrderTypeStop:
		return "Stop"
	}
	return "Unknown"
}
