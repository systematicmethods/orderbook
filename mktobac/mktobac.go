package mktobac

import (
	"github.com/andres-erbsen/clock"
	"orderbook/auction"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/orderbook"
	"orderbook/orderstate"
	"orderbook/tradingevent"
	"time"
)

type obac struct {
	instrument *instrument.Instrument
	ob         orderbook.OrderBook
	ac         auction.OrderBookAuction
	clock      clock.Clock
}

func (b *obac) Instrument() *instrument.Instrument {
	return b.instrument
}

func (b *obac) NewOrder(order *fixmodel.NewOrderSingle) ([]*fixmodel.ExecutionReport, error) {
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
		return execs, err
	} else if execs, err := b.ac.NewOrder(order); err == nil {
		return execs, err
	}
	return nil, nil
}

func (b *obac) CancelOrder(order *fixmodel.OrderCancelRequest) (*fixmodel.ExecutionReport, error) {
	if execs, err := b.ob.CancelOrder(order); err == nil {
		return execs, err
	} else if execs, err := b.ac.CancelOrder(order); err == nil {
		return execs, err
	}
	return nil, nil
}

func (b *obac) Tick(time time.Time) ([]*fixmodel.ExecutionReport, error) {
	return nil, nil
}

func (b *obac) OpenTrading() ([]*fixmodel.ExecutionReport, error) {
	return nil, nil
}

func (b *obac) CloseTrading() (execs []*fixmodel.ExecutionReport, err error) {
	if b.ob.State() == tradingevent.OrderBookStateTradingOpen {
		return b.ob.CloseTrading()
	}
	if b.ac.State() == tradingevent.OrderBookStateAuctionOpen {
		return b.ac.CloseTrading()
	}
	return nil, nil
}

func (b *obac) NoTrading() error {
	return nil
}

func (b *obac) BuyOrders() []*orderstate.OrderState {
	return b.ob.BuyOrders()
}

func (b *obac) SellOrders() []*orderstate.OrderState {
	return b.ob.SellOrders()
}

func (b *obac) BuyAuctionOrders() []*orderstate.OrderState {
	return b.ac.BuyOrders()
}

func (b *obac) SellAuctionOrders() []*orderstate.OrderState {
	return b.ac.SellOrders()
}

func (b *obac) BuySize() int {
	return len(b.ob.BuyOrders())
}

func (b *obac) SellSize() int {
	return len(b.ob.SellOrders())
}

func (b *obac) BuyAuctionSize() int {
	return len(b.ac.BuyOrders())
}

func (b *obac) SellAuctionSize() int {
	return len(b.ac.SellOrders())
}
