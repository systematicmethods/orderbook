package orderbook

type order struct {
	orderevent EventType
	orderid    string
	price      float64
	quantity   int64
	side       Side
	ordertype  OrderType
	data       string
}

type NewOrderEvent interface {
	OrderEvent
}

func MakeNewOrderEvent(orderid string, price float64, ordertype OrderType, side Side, data string) NewOrderEvent {
	return NewOrderEvent(&order{orderevent: EventTypeNewOrder, orderid: orderid, price: price, ordertype: ordertype, side: side, data: data})
}

func (p *order) Orderid() string {
	return p.orderid
}

func (p *order) Price() float64 {
	return p.price
}

func (p *order) Quantity() int64 {
	return p.quantity
}

func (p *order) Data() string {
	return p.data
}

func (p *order) Type() OrderType {
	return p.ordertype
}

func (p *order) Side() Side {
	return p.side
}
