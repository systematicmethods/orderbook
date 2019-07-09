package orderbook

import (
	"github.com/google/uuid"
	"time"
)

type NewOrderSingle interface {
	InstrumentID() string
	ClientID() string
	ClOrdID() string

	Side() Side
	Price() float64
	OrderQty() int64
	OrderType() OrderType
	TimeInForce() TimeInForce
	ExpireOn() time.Time
	TransactTime() time.Time

	OrderID() string

	isBuy() bool
}

type newOrderSingle struct {
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
}

func MakeNewOrderLimit(
	instrumentID string,
	clientID string,
	clOrdID string,
	side Side,
	price float64,
	orderQty int64,
	timeInForce TimeInForce,
	expireOn time.Time,
	transactTime time.Time) NewOrderSingle {
	theOrderID, _ := uuid.NewUUID()
	return NewOrderSingle(&newOrderSingle{
		instrumentID,
		clientID,
		clOrdID,
		side,
		price,
		orderQty,
		OrderTypeLimit,
		timeInForce,
		expireOn,
		transactTime,
		EventTypeNewOrderSingle,
		theOrderID.String(),
	})
}

func MakeNewOrderMarket(
	instrumentID string,
	clientID string,
	clOrdID string,
	side Side,
	orderQty int64,
	timeInForce TimeInForce,
	expireOn time.Time,
	transactTime time.Time) NewOrderSingle {
	theOrderID, _ := uuid.NewUUID()
	return NewOrderSingle(&newOrderSingle{
		instrumentID,
		clientID,
		clOrdID,
		side,
		0,
		orderQty,
		OrderTypeMarket,
		timeInForce,
		expireOn,
		transactTime,
		EventTypeNewOrderSingle,
		theOrderID.String(),
	})
}

func (p *newOrderSingle) OrderID() string {
	return p.orderID
}

func (p *newOrderSingle) Price() float64 {
	return p.price
}

func (p *newOrderSingle) OrderQty() int64 {
	return p.orderQty
}

func (p *newOrderSingle) OrderType() OrderType {
	return p.orderType
}

func (p *newOrderSingle) Side() Side {
	return p.side
}

func (p *newOrderSingle) ClOrdID() string {
	return p.clOrdID
}

func (p *newOrderSingle) ClientID() string {
	return p.clientID
}

func (p *newOrderSingle) InstrumentID() string {
	return p.instrumentID
}

func (p *newOrderSingle) TimeInForce() TimeInForce {
	return p.timeInForce
}

func (p *newOrderSingle) ExpireOn() time.Time {
	return p.expireOn
}

func (p *newOrderSingle) TransactTime() time.Time {
	return p.transactTime
}

func (p *newOrderSingle) isBuy() bool {
	return p.Side() == SideBuy
}
