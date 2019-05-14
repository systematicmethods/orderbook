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
	ClOrdID() string
	OrderID() string
	Price() float64
	Timestamp() time.Time
	UUID() uuid.UUID
	Data() string
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

func (p *order) ClOrdID() string {
	return p.clOrdID
}

func (p *order) OrderID() string {
	return p.orderID
}

func (p *order) Price() float64 {
	return p.price
}

func (p *order) Timestamp() time.Time {
	return time.Unix(p.timestamp.Time().UnixTime())
}

func (p *order) UUID() uuid.UUID {
	return p.timestamp
}

func (p *order) Data() string {
	return p.data
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
