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

	matchOrder() []ExecutionReport
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
		//fmt.Printf("NewOrder added %v added qty %d\n", order, order.OrderQty())
		execs = append(execs, MakeNewOrderAckExecutionReport(order))
		//fmt.Printf("NewOrder execs in OrderBook %v\n", execs)
		if order.Side() == SideBuy {
			b.buyOrders.Add(order)
		} else {
			b.sellOrders.Add(order)
		}
		execs = append(execs, b.matchOrder()...)
		//fmt.Printf("NewOrder execs after in OrderBook %v", execs)
		return execs, nil
	}
	return nil, nil
}

func (b *orderbook) matchOrder() []ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", b.sellOrders.Size(), b.buyOrders.Size())
	matchexecs := []ExecutionReport{}
	for buyiter := b.buyOrders.iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(OrderState)
		if buyorder.Side() == SideBuy && b.sellOrders.Size() > 0 {
			//fmt.Printf("buy order %s %f %d\n", SideToString(buyorder.Side()), buyorder.Price(), buyorder.LeavesQty())
			for selliter := b.sellOrders.iterator(); selliter.Next() == true; {
				sellorder := selliter.Value().(OrderState)
				//fmt.Printf("buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
				if buyorder.Price() >= sellorder.Price() {
					toFill := min(sellorder.LeavesQty(), buyorder.LeavesQty())
					if buyorder.fill(toFill) {
						b.buyOrders.RemoveByID(buyorder.OrderID())
					}
					matchexecs = append(matchexecs, MakeFillExecutionReport(buyorder, sellorder.Price(), toFill))
					if sellorder.fill(toFill) {
						b.sellOrders.RemoveByID(sellorder.OrderID())
					}
					matchexecs = append(matchexecs, MakeFillExecutionReport(sellorder, sellorder.Price(), toFill))
				}
				//fmt.Printf("After loop buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
			}
		}
	}
	return matchexecs
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

func min(x int64, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
