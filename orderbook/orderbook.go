package orderbook

import (
	"github.com/google/uuid"
	"orderbook/instrument"
	"orderbook/orderlist"
)

type OrderBook interface {
	Instrument() *instrument.Instrument
	NewOrder(order NewOrderEvent) error
	BuySize() int
	SellSize() int
	BuyOrders() []orderlist.Order
	SellOrders() []orderlist.Order
}

func MakeOrderBook(instrument instrument.Instrument) OrderBook {
	b := orderbook{instrument: &instrument}
	b.buyOrders = orderlist.NewOrderList(orderlist.HighToLow)
	b.sellOrders = orderlist.NewOrderList(orderlist.HighToLow)
	return OrderBook(&b)
}

type orderbook struct {
	instrument *instrument.Instrument
	buyOrders  orderlist.OrderList
	sellOrders orderlist.OrderList
}

func (b *orderbook) NewOrder(neworder NewOrderEvent) error {
	if neworder.Orderid() != "" {
		order := orderlist.NewOrder(neworder.Orderid(), neworder.Price(), newID(uuid.NewUUID()), "data")
		if neworder.Side() == SideBuy {
			b.buyOrders.Add(order)
		} else {
			b.sellOrders.Add(order)
		}
	}
	return nil
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

func (b *orderbook) BuyOrders() []orderlist.Order {
	return b.buyOrders.Orders()
}

func (b *orderbook) SellOrders() []orderlist.Order {
	return b.sellOrders.Orders()
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}
