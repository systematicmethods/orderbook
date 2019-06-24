package orderbook

import (
	"github.com/google/uuid"
	"orderbook/instrument"
	"time"
)

type OrderBook interface {
	Instrument() *instrument.Instrument
	NewOrder(order NewOrderSingle) ([]ExecutionReport, error)
	CancelOrder(order OrderCancelRequest) (ExecutionReport, error)
	BuySize() int
	SellSize() int
	BuyOrders() []OrderState
	SellOrders() []OrderState
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

func (b *orderbook) NewOrder(order NewOrderSingle) ([]ExecutionReport, error) {
	if order.OrderID() != "" {
		execs := []ExecutionReport{}
		order := NewOrder(order, newID(uuid.NewUUID()), time.Now())
		if order.Side() == SideBuy {
			b.buyOrders.Add(order)
		} else {
			b.sellOrders.Add(order)
		}
		execs = append(execs, MakeNewOrderAckExecutionReport(order))
		//if order.Side() == SideBuy {
		//	b.buyOrders.Add(order)
		//} else {
		//	b.sellOrders.Add(order)
		//}

		return execs, nil
	}
	return nil, nil
}

func (b *orderbook) CancelOrder(order OrderCancelRequest) (ExecutionReport, error) {
	var ord OrderState
	if order.Side() == SideBuy {
		ord = b.buyOrders.FindByClOrdID(order.OrigClOrdID())
		if ord != nil {
			b.buyOrders.RemoveByID(ord.OrderID())
			return MakeCancelOrderExecutionReport(ord, order), nil
		}
	} else {
		ord = b.sellOrders.FindByClOrdID(order.OrigClOrdID())
		if ord != nil {
			b.sellOrders.RemoveByID(ord.OrderID())
			return MakeCancelOrderExecutionReport(ord, order), nil
		}
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

func (b *orderbook) BuyOrders() []OrderState {
	return b.buyOrders.Orders()
}

func (b *orderbook) SellOrders() []OrderState {
	return b.sellOrders.Orders()
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}
