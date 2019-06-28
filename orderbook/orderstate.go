package orderbook

import (
	"fmt"
	"github.com/google/uuid"
	"orderbook/uuidext"
	"time"
)

type OrderState interface {
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
	CreatedOn() time.Time
	UpdatedOn() time.Time

	OrderID() string
	Timestamp() uuid.UUID

	LeavesQty() int64
	CumQty() int64
	OrdStatus() OrdStatus
}

func newOrderForTesting(clOrdID string, orderID string, price float64, timestamp uuid.UUID, data string) OrderState {
	return &orderState{clOrdID: clOrdID, orderID: orderID, price: price, data: data, timestamp: timestamp}
}

func NewOrder(ord NewOrderSingle, timestamp uuid.UUID, createdOn time.Time) OrderState {
	return &orderState{
		ord.InstrumentID(),
		ord.ClientID(),
		ord.ClOrdID(),
		ord.Side(),
		ord.Price(),
		ord.OrderQty(),
		ord.OrderType(),
		ord.TimeInForce(),
		ord.ExpireOn(),
		ord.TransactTime(),
		createdOn,
		time.Time{},
		ord.OrderID(),
		timestamp,
		ord.OrderQty(),
		0,
		OrdStatusNew,
		"",
	}
}

type orderState struct {
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

	createdOn time.Time
	updatedOn time.Time

	orderID   string
	timestamp uuid.UUID

	leavesQty int64
	cumQty    int64
	ordStatus OrdStatus
	data      string
}

func (b *orderState) String() string {
	str1 := fmt.Sprintf("OrderState: instrumentID:%s, clientID:%s, clOrdID:%s, side:%v, pricee:%f, orderQty:%d, orderType:%v, timeInForce:%v, expireOn:%v, transactTime:%v, ",
		b.instrumentID,
		b.clientID,
		b.clOrdID,
		SideToString(b.side),
		b.price,
		b.orderQty,
		OrderTypeToString(b.orderType),
		b.timeInForce,
		b.expireOn,
		b.transactTime)
	str2 := fmt.Sprintf("createdOn:%v, updatedOn:%v, orderID:%s, timestamp:%v, leavesQty:%d, cumQty:%d, ordStatus:%v",
		b.createdOn,
		b.updatedOn,
		b.orderID,
		b.timestamp,
		b.leavesQty,
		b.cumQty,
		OrdStatusToString(b.ordStatus))
	return fmt.Sprintf("%s %s", str1, str2)
}

func (p *orderState) InstrumentID() string {
	return p.instrumentID
}

func (p *orderState) ClientID() string {
	return p.clientID
}

func (p *orderState) ClOrdID() string {
	return p.clOrdID
}

func (p *orderState) Side() Side {
	return p.side
}

func (p *orderState) Price() float64 {
	return p.price
}

func (p *orderState) OrderQty() int64 {
	return p.orderQty
}

func (p *orderState) OrderType() OrderType {
	return p.orderType
}

func (p *orderState) TimeInForce() TimeInForce {
	return p.timeInForce
}

func (p *orderState) ExpireOn() time.Time {
	return p.expireOn
}

func (p *orderState) TransactTime() time.Time {
	return p.transactTime
}

func (p *orderState) CreatedOn() time.Time {
	return p.createdOn
}

func (p *orderState) UpdatedOn() time.Time {
	return p.updatedOn
}

func (p *orderState) OrderID() string {
	return p.orderID
}

func (p *orderState) Timestamp() uuid.UUID {
	return p.timestamp
}

func (p *orderState) LeavesQty() int64 {
	return p.leavesQty
}

func (p *orderState) CumQty() int64 {
	return p.cumQty
}

func (p *orderState) OrdStatus() OrdStatus {
	return p.ordStatus
}

func sellPriceComparator(a, b interface{}) int {
	apti := a.(*orderState)
	bpti := b.(*orderState)
	switch {
	case apti.price > bpti.price:
		return 1
	case apti.price < bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}

func buyPriceComparator(a, b interface{}) int {
	apti := a.(*orderState)
	bpti := b.(*orderState)
	switch {
	case apti.price < bpti.price:
		return 1
	case apti.price > bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}
