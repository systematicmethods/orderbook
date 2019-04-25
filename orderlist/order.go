package orderlist

import (
	"github.com/google/uuid"
	"orderbook/uuidext"
	"time"
)

type order struct {
	orderid  string
	price    float64
	timeuuid uuid.UUID
	data     string
}

type Order interface {
	Orderid() string
	Price() float64
	Timestamp() time.Time
	Data() string
}

func NewOrder(orderid string, price float64, timestamp uuid.UUID, data string) *order {
	return &order{orderid: orderid, price: price, data: data, timeuuid: timestamp}
}

func (p *order) Orderid() string {
	return p.orderid
}

func (p *order) Price() float64 {
	return p.price
}

func (p *order) Timestamp() time.Time {
	return time.Unix(p.timeuuid.Time().UnixTime())
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
