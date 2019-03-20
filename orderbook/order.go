package orderbook

import (
	"github.com/google/uuid"
	"time"
)

type order struct {
	orderid  string
	price    float64
	timmuuid uuid.UUID
	data     string
}

type Order interface {
	Orderid() string
	Price() float64
	Timestamp() time.Time
	Data() string
}

func NewOrder(orderid string, price float64, timestamp time.Time, data string) *order {
	uuid, _ := uuid.NewUUID()
	return &order{orderid: orderid, price: price, data: data, timmuuid: uuid}
}

func (p *order) Orderid() string {
	return p.orderid
}

func (p *order) Price() float64 {
	return p.price
}

func (p *order) Timestamp() time.Time {
	return time.Unix(p.timmuuid.Time().UnixTime())
	//time := p.timmuuid.Time()
	//return time
}

func (p *order) Data() string {
	return p.data
}

func priceComparator(a, b interface{}) int {
	apti := a.(*order)
	bpti := b.(*order)
	switch {
	case apti.price > bpti.price:
		return 1
	case apti.price < bpti.price:
		return -1
	default:
		switch {
		case apti.timmuuid.Time() > bpti.timmuuid.Time(): // after
			return 1
		case apti.timmuuid.Time() < bpti.timmuuid.Time(): // before
			return -1
		default:
			return 0
		}
	}
}
