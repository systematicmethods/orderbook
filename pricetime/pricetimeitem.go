package pricetime

import (
	"github.com/google/uuid"
	"time"
)

type priceTimeItem struct {
	orderid   string
	price     float64
	timestamp time.Time
	timmuuid  uuid.UUID
	data      string
}

type PriceTimeItem interface {
	Orderid() string
	Price() float64
	Timestamp() time.Time
	Data() string
}

func NewPriceTimeItem(orderid string, price float64, timestamp time.Time, data string) *priceTimeItem {
	uuid, _ := uuid.NewUUID()
	return &priceTimeItem{orderid: orderid, price: price, data: data, timestamp: timestamp, timmuuid: uuid}
}

func (p *priceTimeItem) Orderid() string {
	return p.orderid
}

func (p *priceTimeItem) Price() float64 {
	return p.price
}

func (p *priceTimeItem) Timestamp() time.Time {
	return p.timestamp
}

func (p *priceTimeItem) Data() string {
	return p.data
}

func priceComparator(a, b interface{}) int {
	apti := a.(*priceTimeItem)
	bpti := b.(*priceTimeItem)
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
