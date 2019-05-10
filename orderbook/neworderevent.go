package orderbook

type NewOrderEvent interface {
	OrderEvent
}

func MakeNewOrderEvent(orderid string, price float64, ordertype OrderType, side Side, data string) NewOrderEvent {
	return NewOrderEvent(&orderEvent{eventType: EventTypeNewOrder, orderID: orderid, price: price, orderType: ordertype, side: side, data: data})
}
