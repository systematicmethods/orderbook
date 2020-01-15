package mktstopob

import (
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"orderbook/etype"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/orderbook"
	"orderbook/orderstate"
	"orderbook/tradingevent"
	"time"
)

type stopob struct {
	instrument *instrument.Instrument
	ob         orderbook.OrderBook
	so         BuySellOrders

	clock clock.Clock
}

type BuySellOrders struct {
	BuyOrders  *orderstate.Orderlist
	SellOrders *orderstate.Orderlist
}

const (
	rejected = etype.Error("rejected")
)

func (b *stopob) Instrument() *instrument.Instrument {
	return b.instrument
}

func (b *stopob) NewOrder(order *fixmodel.NewOrderSingle) ([]*fixmodel.ExecutionReport, error) {
	execs := []*fixmodel.ExecutionReport{}

	if b.ob.State() == tradingevent.OrderBookStateTradingClosed {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, "Trading closed"))
		return execs, rejected
	}

	if order.OrderType() == fixmodel.OrderTypeStop {
		return addNewOrder(order, &b.so, b.clock)
	}
	//if b.ob.State() == tradingevent.OrderBookStateTradingClosed {
	//	execs := append([]*fixmodel.ExecutionReport{}, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, ""))
	//	return execs, nil
	//}
	//if b.ob. == tradingevent.OrderBookStateTradingClosed && order.TimeInForce() != fixmodel.TimeInForceGoodForAuction {
	//}
	//
	//if b.orderBookState == tradingevent.OrderBookStateAuctionClosed && order.TimeInForce() == fixmodel.TimeInForceGoodForAuction {
	//	execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonDuplicateOrder, ""))
	//	return execs, nil
	//}
	//
	//if order.TimeInForce() == fixmodel.TimeInForceGoodForAuction {
	//	return addNewOrder(order, &b.auctionOrders, b.clock)
	//}

	if execs, err := b.ob.NewOrder(order); err == nil {
		return b.triggerStop(execs), err
	}
	return nil, nil
}

func (b *stopob) CancelOrder(order *fixmodel.OrderCancelRequest) (*fixmodel.ExecutionReport, error) {
	if execs, err := b.ob.CancelOrder(order); err == nil {
		return execs, err
	}
	return nil, nil
}

func (b *stopob) Tick(time time.Time) ([]*fixmodel.ExecutionReport, error) {
	return nil, nil
}

func (b *stopob) OpenTrading() ([]*fixmodel.ExecutionReport, error) {
	return nil, nil
}

func (b *stopob) CloseTrading() (execs []*fixmodel.ExecutionReport, err error) {
	if b.ob.State() == tradingevent.OrderBookStateTradingOpen {
		return b.ob.CloseTrading()
	}
	return nil, nil
}

func (b *stopob) NoTrading() error {
	return nil
}

func (b *stopob) BuyOrders() []*orderstate.OrderState {
	return b.ob.BuyOrders()
}

func (b *stopob) SellOrders() []*orderstate.OrderState {
	return b.ob.SellOrders()
}

func (b *stopob) BuyStopOrders() []*orderstate.OrderState {
	return b.so.BuyOrders.Orders()
}

func (b *stopob) SellStopOrders() []*orderstate.OrderState {
	return b.so.SellOrders.Orders()
}

func (b *stopob) BuySize() int {
	return len(b.ob.BuyOrders())
}

func (b *stopob) SellSize() int {
	return len(b.ob.SellOrders())
}

func (b *stopob) BuyStopSize() int {
	return b.so.BuyOrders.Size()
}

func (b *stopob) SellStopSize() int {
	return b.so.SellOrders.Size()
}

func (b *stopob) triggerStop(execs []*fixmodel.ExecutionReport) []*fixmodel.ExecutionReport {
	if b.ob.BuySize() > 0 {
		executedstoporders := []*orderstate.OrderState{}
		for iter := b.so.SellOrders.Iterator(); iter.Next() == true; {
			order := iter.Value().(*orderstate.OrderState)
			// submit sell order if buy top of book is less than sell stop order price
			if b.ob.BuyTop().Price < order.Price() {
				//if order.Price() < b.ob.BuyOrders
				neworder := newNeworder(order, b.clock)
				aexecs, _ := b.ob.NewOrder(neworder)
				execs = append(execs, aexecs...)
				executedstoporders = append(executedstoporders, order)

			}
		}
		for _, v := range executedstoporders {
			b.so.SellOrders.RemoveByID(v.OrderID())
		}
	}
	if b.ob.SellSize() > 0 {
		executedstoporders := []*orderstate.OrderState{}
		for iter := b.so.BuyOrders.Iterator(); iter.Next() == true; {
			order := iter.Value().(*orderstate.OrderState)
			// submit buy order if sell top of book is greater than buy stop order price
			if b.ob.SellTop().Price > order.Price() {
				neworder := newNeworder(order, b.clock)
				aexecs, _ := b.ob.NewOrder(neworder)
				execs = append(execs, aexecs...)
				executedstoporders = append(executedstoporders, order)
			}
		}
		for _, v := range executedstoporders {
			b.so.BuyOrders.RemoveByID(v.OrderID())
		}
	}
	//if  b.ob.SellSize() > 0 {
	//}
	return execs
}

func newNeworder(order *orderstate.OrderState, clock clock.Clock) *fixmodel.NewOrderSingle {
	return fixmodel.NewNewOrder(order.InstrumentID(),
		order.ClientID(),
		order.ClOrdID(),
		order.Side(),
		order.Price(),
		order.OrderQty(),
		order.TimeInForce(),
		order.ExpireOn(),
		clock.Now(),
		fixmodel.OrderTypeMarket)
}

func addNewOrder(order *fixmodel.NewOrderSingle, bs *BuySellOrders, clock clock.Clock) ([]*fixmodel.ExecutionReport, error) {
	execs := []*fixmodel.ExecutionReport{}
	var err error
	neworder := orderstate.NewOrder(order, newID(uuid.NewUUID()), clock.Now())
	execs = append(execs, orderstate.NewNewOrderAckExecutionReport(neworder))
	if order.Side() == fixmodel.SideBuy {
		err = bs.BuyOrders.Add(neworder)
	} else {
		err = bs.SellOrders.Add(neworder)
	}
	return execs, err
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}
