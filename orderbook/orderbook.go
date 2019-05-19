package orderbook

import (
	"github.com/google/uuid"
	"orderbook/instrument"
	"time"
)

type OrderBook interface {
	Instrument() *instrument.Instrument
	NewOrder(order NewOrderSingle) (ExecutionReport, error)
	BuySize() int
	SellSize() int
	BuyOrders() []Order
	SellOrders() []Order
}

func MakeOrderBook(instrument instrument.Instrument) OrderBook {
	b := orderbook{instrument: &instrument}
	b.buyOrders = NewOrderList(HighToLow)
	b.sellOrders = NewOrderList(HighToLow)
	return OrderBook(&b)
}

type orderbook struct {
	instrument *instrument.Instrument
	buyOrders  OrderList
	sellOrders OrderList
}

func (b *orderbook) NewOrder(neworder NewOrderSingle) (ExecutionReport, error) {
	if neworder.OrderID() != "" {
		order := NewOrder(neworder, newID(uuid.NewUUID()), time.Now())
		if neworder.Side() == SideBuy {
			b.buyOrders.Add(order)
		} else {
			b.sellOrders.Add(order)
		}
		return MakeNewOrderAckExecutionReport(order), nil
	}
	return nil, nil
}

func (b *orderbook) BuySize() int {
	return b.buyOrders.Size()
}

func (b *orderbook) SellSize() int {
	return b.sellOrders.Size()
}

func (b *orderbook) Instrument() *instrument.Instrument {
	return b.instrument
}

func (b *orderbook) BuyOrders() []Order {
	return b.buyOrders.Orders()
}

func (b *orderbook) SellOrders() []Order {
	return b.sellOrders.Orders()
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}
