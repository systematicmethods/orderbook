package orderbook

import (
	"github.com/google/uuid"
	"orderbook/uuidext"
	"time"
)

type order struct {
	timeuuid     uuid.UUID
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
	Orderid() string
	Price() float64
	Timestamp() time.Time
	UUID() uuid.UUID
	Data() string
}

func NewOrder(orderID string, price float64, timestamp uuid.UUID, data string) *order {
	return &order{orderID: orderID, price: price, data: data, timeuuid: timestamp}
}

func (p *order) Orderid() string {
	return p.orderID
}

func (p *order) Price() float64 {
	return p.price
}

func (p *order) Timestamp() time.Time {
	return time.Unix(p.timeuuid.Time().UnixTime())
}

func (p *order) UUID() uuid.UUID {
	return p.timeuuid
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
		return uuidext.UUIDComparator(apti.timeuuid, bpti.timeuuid)
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
		return uuidext.UUIDComparator(apti.timeuuid, bpti.timeuuid)
	}
}
