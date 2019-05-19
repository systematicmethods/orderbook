package orderbook

import (
	"github.com/google/uuid"
	"orderbook/uuidext"
	"time"
)

type order struct {
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

type Order interface {
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

func NewOrder2(clOrdID string, orderID string, price float64, timestamp uuid.UUID, data string) Order {
	return &order{clOrdID: clOrdID, orderID: orderID, price: price, data: data, timestamp: timestamp}
}

func NewOrder(ord NewOrderSingle, timestamp uuid.UUID, createdOn time.Time) Order {
	return &order{
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

func (p *order) InstrumentID() string {
	return p.instrumentID
}

func (p *order) ClientID() string {
	return p.clientID
}

func (p *order) ClOrdID() string {
	return p.clOrdID
}

func (p *order) Side() Side {
	return p.side
}

func (p *order) Price() float64 {
	return p.price
}

func (p *order) OrderQty() int64 {
	return p.orderQty
}

func (p *order) OrderType() OrderType {
	return p.orderType
}

func (p *order) TimeInForce() TimeInForce {
	return p.timeInForce
}

func (p *order) ExpireOn() time.Time {
	return p.expireOn
}

func (p *order) TransactTime() time.Time {
	return p.transactTime
}

func (p *order) CreatedOn() time.Time {
	return p.createdOn
}

func (p *order) UpdatedOn() time.Time {
	return p.updatedOn
}

func (p *order) OrderID() string {
	return p.orderID
}

func (p *order) Timestamp() uuid.UUID {
	return p.timestamp
}

func (p *order) LeavesQty() int64 {
	return p.leavesQty
}

func (p *order) CumQty() int64 {
	return p.cumQty
}

func (p *order) OrdStatus() OrdStatus {
	return p.ordStatus
}

func sellPriceComparator(a, b interface{}) int {
	apti := a.(*order)
	bpti := b.(*order)
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
	apti := a.(*order)
	bpti := b.(*order)
	switch {
	case apti.price < bpti.price:
		return 1
	case apti.price > bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}
