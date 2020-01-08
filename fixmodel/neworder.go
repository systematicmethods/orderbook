package fixmodel

import (
	"github.com/google/uuid"
	"time"
)

type NewOrderSingle struct {
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

func NewNewOrder(
	instrumentID string,
	clientID string,
	clOrdID string,
	side Side,
	price float64,
	orderQty int64,
	timeInForce TimeInForce,
	expireOn time.Time,
	transactTime time.Time,
	orderType OrderType) *NewOrderSingle {
	theOrderID, _ := uuid.NewUUID()
	return &NewOrderSingle{
		instrumentID,
		clientID,
		clOrdID,
		side,
		price,
		orderQty,
		orderType,
		timeInForce,
		expireOn,
		transactTime,
		EventTypeNewOrderSingle,
		theOrderID.String(),
	}
}

func (p *NewOrderSingle) OrderID() string {
	return p.orderID
}

func (p *NewOrderSingle) Price() float64 {
	return p.price
}

func (p *NewOrderSingle) OrderQty() int64 {
	return p.orderQty
}

func (p *NewOrderSingle) OrderType() OrderType {
	return p.orderType
}

func (p *NewOrderSingle) Side() Side {
	return p.side
}

func (p *NewOrderSingle) ClOrdID() string {
	return p.clOrdID
}

func (p *NewOrderSingle) ClientID() string {
	return p.clientID
}

func (p *NewOrderSingle) InstrumentID() string {
	return p.instrumentID
}

func (p *NewOrderSingle) TimeInForce() TimeInForce {
	return p.timeInForce
}

func (p *NewOrderSingle) ExpireOn() time.Time {
	return p.expireOn
}

func (p *NewOrderSingle) TransactTime() time.Time {
	return p.transactTime
}

func (p *NewOrderSingle) isBuy() bool {
	return p.Side() == SideBuy
}
