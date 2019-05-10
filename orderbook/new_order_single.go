package orderbook

import "time"

type NewOrderSingle interface {
	OrderID() string
	Price() float64
	OrderQty() int64
	OrderType() OrderType
	Side() Side
	ClOrdID() string
	ClientID() string
	InstrumentID() string
	TimeInForce() TimeInForce
	ExpireOn() time.Time
	TransactTime() time.Time
	Data() string
}

func MakeNewOrderEvent(orderid string, price float64, ordertype OrderType, side Side, data string) NewOrderSingle {
	return NewOrderSingle(&orderEvent{eventType: EventTypeNewOrder, orderID: orderid, price: price, orderType: ordertype, side: side, data: data})
}

type orderEvent struct {
	instrumentID string
	clientID     string
	clOrdID      string
	side         Side

	price        float64
	orderQty     int64
	orderType    OrderType
	timeInForce  TimeInForce
	expireOn     time.Time
	transactTime time.Time

	eventType EventType
	orderID   string
	data      string
}

func (p *orderEvent) OrderID() string {
	return p.orderID
}

func (p *orderEvent) Price() float64 {
	return p.price
}

func (p *orderEvent) OrderQty() int64 {
	return p.orderQty
}

func (p *orderEvent) Data() string {
	return p.data
}

func (p *orderEvent) OrderType() OrderType {
	return p.orderType
}

func (p *orderEvent) Side() Side {
	return p.side
}

func (p *orderEvent) ClOrdID() string {
	return p.clOrdID
}

func (p *orderEvent) ClientID() string {
	return p.clientID
}

func (p *orderEvent) InstrumentID() string {
	return p.instrumentID
}

func (p *orderEvent) TimeInForce() TimeInForce {
	return p.timeInForce
}

func (p *orderEvent) ExpireOn() time.Time {
	return p.expireOn
}

func (p *orderEvent) TransactTime() time.Time {
	return p.transactTime
}
