package orderbook

import "time"

type OrderEvent interface {
	OrderID() string
	Price() float64
	Quantity() int64
	OrderType() OrderType
	Side() Side
	ClientOrderID() string
	ClientID() string
	InstrumentID() string
	TimeInForce() TimeInForce
	ExpireOn() time.Time
	TransactTime() time.Time
	Data() string
}

type orderEvent struct {
	eventType     EventType
	orderID       string
	price         float64
	quantity      int64
	side          Side
	orderType     OrderType
	clientOrderID string
	clientID      string
	instrumentID  string
	timeInForce   TimeInForce
	expireOn      time.Time
	transactTime  time.Time
	data          string
}

func (p *orderEvent) OrderID() string {
	return p.orderID
}

func (p *orderEvent) Price() float64 {
	return p.price
}

func (p *orderEvent) Quantity() int64 {
	return p.quantity
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

func (p *orderEvent) ClientOrderID() string {
	return p.clientOrderID
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
