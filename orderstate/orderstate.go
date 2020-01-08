package orderstate

import (
	"fmt"
	"github.com/google/uuid"
	"orderbook/fixmodel"
	"orderbook/uuidext"
	"time"
)

type OrderState struct {
	instrumentID string
	clientID     string
	clOrdID      string

	side         fixmodel.Side
	price        float64
	orderQty     int64
	orderType    fixmodel.OrderType
	timeInForce  fixmodel.TimeInForce
	expireOn     time.Time
	transactTime time.Time

	createdOn time.Time
	updatedOn time.Time

	orderID   string
	timestamp uuid.UUID

	leavesQty int64
	cumQty    int64
	ordStatus fixmodel.OrdStatus
	data      string
}

func NewOrderForTesting(clOrdID string, orderID string, price float64, timestamp uuid.UUID) *OrderState {
	return &OrderState{clOrdID: clOrdID, orderID: orderID, price: price, data: "", timestamp: timestamp}
}

func NewOrder(ord *fixmodel.NewOrderSingle, timestamp uuid.UUID, createdOn time.Time) *OrderState {
	return &OrderState{
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
		fixmodel.OrdStatusNew,
		"",
	}
}

func NewOrderState(
	instrumentID string,
	clientID string,
	clOrdID string,
	side fixmodel.Side,
	price float64,
	orderQty int64,
	orderType fixmodel.OrderType,
	timeInForce fixmodel.TimeInForce,
	expireOn time.Time,
	transactTime time.Time,
	createdOn time.Time,
	updatedOn time.Time,
	orderID string,
	timestamp uuid.UUID,
	leavesQty int64,
	cumQty int64,
	ordStatus fixmodel.OrdStatus,
) *OrderState {
	return &OrderState{
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
		createdOn,
		updatedOn,
		orderID,
		timestamp,
		leavesQty,
		cumQty,
		ordStatus,
		"",
	}
}

func (b *OrderState) String() string {
	if b == nil {
		return "OrderState is nil"
	}
	str1 := fmt.Sprintf("OrderState: instrumentID:%s, clientID:%s, clOrdID:%s, side:%v, price:%f, orderQty:%d, orderType:%v, timeInForce:%v, expireOn:%v, transactTime:%v, ",
		b.instrumentID,
		b.clientID,
		b.clOrdID,
		b.side,
		b.price,
		b.orderQty,
		b.orderType,
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
		fixmodel.OrdStatusToString(b.ordStatus))
	return fmt.Sprintf("%s %s", str1, str2)
}

func (p *OrderState) InstrumentID() string {
	return p.instrumentID
}

func (p *OrderState) ClientID() string {
	return p.clientID
}

func (p *OrderState) ClOrdID() string {
	return p.clOrdID
}

func (p *OrderState) Side() fixmodel.Side {
	return p.side
}

func (p *OrderState) Price() float64 {
	return p.price
}

func (p *OrderState) OrderQty() int64 {
	return p.orderQty
}

func (p *OrderState) OrderType() fixmodel.OrderType {
	return p.orderType
}

func (p *OrderState) TimeInForce() fixmodel.TimeInForce {
	return p.timeInForce
}

func (p *OrderState) ExpireOn() time.Time {
	return p.expireOn
}

func (p *OrderState) TransactTime() time.Time {
	return p.transactTime
}

func (p *OrderState) CreatedOn() time.Time {
	return p.createdOn
}

func (p *OrderState) UpdatedOn() time.Time {
	return p.updatedOn
}

func (p *OrderState) OrderID() string {
	return p.orderID
}

func (p *OrderState) Timestamp() uuid.UUID {
	return p.timestamp
}

func (p *OrderState) LeavesQty() int64 {
	return p.leavesQty
}

func (p *OrderState) CumQty() int64 {
	return p.cumQty
}

func (p *OrderState) OrdStatus() fixmodel.OrdStatus {
	return p.ordStatus
}

func (o *OrderState) Fill(qty int64) bool {
	//fmt.Printf("Fill OrderState %v qty %d\n", o, qty)
	o.cumQty = o.cumQty + qty
	o.leavesQty = o.leavesQty - qty
	if o.leavesQty <= 0 {
		o.ordStatus = fixmodel.OrdStatusFilled
		return true
	} else {
		o.ordStatus = fixmodel.OrdStatusPartiallyFilled
	}
	return false
}

func SellPriceComparator(a, b interface{}) int {
	apti := a.(*OrderState)
	bpti := b.(*OrderState)
	switch {
	case apti.price > bpti.price:
		return 1
	case apti.price < bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}

func BuyPriceComparator(a, b interface{}) int {
	apti := a.(*OrderState)
	bpti := b.(*OrderState)
	switch {
	case apti.price < bpti.price:
		return 1
	case apti.price > bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}
