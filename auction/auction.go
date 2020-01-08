package auction

import (
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/orderbook"
	"orderbook/orderstate"
	"orderbook/tradingevent"
	"time"
)

type OrderBookAuction interface {
	orderbook.OrderBook
}

type BuySellOrders struct {
	BuyOrders  *orderstate.Orderlist
	SellOrders *orderstate.Orderlist
}

type auction struct {
	instrument     *instrument.Instrument
	orders         BuySellOrders
	orderBookState tradingevent.OrderBookState
	clock          clock.Clock
}

func (b *auction) Instrument() *instrument.Instrument {
	return b.instrument
}

func (b *auction) State() tradingevent.OrderBookState {
	return b.orderBookState
}

func (b *auction) NewOrder(order *fixmodel.NewOrderSingle) ([]*fixmodel.ExecutionReport, error) {
	//if order.TimeInForce() != fixmodel.TimeInForceGoodForAuction {
	//	return nil, nil
	//}

	execs := []*fixmodel.ExecutionReport{}

	if b.orderBookState == tradingevent.OrderBookStateOrderEntryClosed {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, "Exchange closed"))
		return execs, nil
	}
	if b.orderBookState != tradingevent.OrderBookStateAuctionOpen {
		execs := append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, "Auction not open"))
		return execs, nil
	} else if order.TimeInForce() == fixmodel.TimeInForceGoodForAuction {
		return addNewOrder(order, &b.orders, b.clock)
	} else {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonUnsupportedOrderCharacteristic, "Must be an auction order"))
		return execs, nil
	}
}

func (b *auction) CancelOrder(order *fixmodel.OrderCancelRequest) (*fixmodel.ExecutionReport, error) {
	return nil, nil
}

func (b *auction) Tick(time time.Time) ([]*fixmodel.ExecutionReport, error) {
	return nil, nil
}

func (b *auction) OpenTrading() ([]*fixmodel.ExecutionReport, error) {
	var err error
	b.orderBookState, err = tradingevent.OrderBookStateChange(b.orderBookState, tradingevent.OrderBookEventTypeOpenAuction)
	return nil, err
}

func (b *auction) CloseTrading() (execs []*fixmodel.ExecutionReport, err error) {
	execs, _, _, err = b.CloseAuction()
	return
}

func (b *auction) NoTrading() error {
	return nil
}

/*
	1	Find minimum price on buy side that match on the sell side
	2	Find maximum price on sell side that match on the buy side
	3	Find max volume that match between min and max price
	4	Match orders in price range and max volume on buy and sell side - use price time priority when orders are same price
	5	Calculate vwap buy orders using max volume
	6	Calculate vwap sell orders using max volume
	7	(vwap buy + vwap sell) /2
	8	Round to 2 decimals
	9	Fill orders to rounded auction price
	10	Cancel remaining orders
*/
func (b *auction) CloseAuction() (execs []*fixmodel.ExecutionReport, clearingPrice float64, clearingVol int64, err error) {
	execs = []*fixmodel.ExecutionReport{}
	state := newAuctionCloseCalculator()
	b.orderBookState, err = tradingevent.OrderBookStateChange(b.orderBookState, tradingevent.OrderBookEventTypeCloseAuction)
	if err == nil {
		var exs []*fixmodel.ExecutionReport
		exs, err = state.fillAuctionAtClearingPrice(&b.orders)
		execs = append(execs, exs...)
		clearingPrice, _ = state.state().midclearingprice.Float64()
		clearingVol = state.state().clearingvol
		exs = cancelOrders(&b.orders)
		execs = append(execs, exs...)
	}
	return
}

func (b *auction) BuyOrders() []*orderstate.OrderState {
	return b.orders.BuyOrders.Orders()
}

func (b *auction) SellOrders() []*orderstate.OrderState {
	return b.orders.SellOrders.Orders()
}

func (b *auction) SellSize() int {
	return b.orders.SellOrders.Size()
}

func (b *auction) BuySize() int {
	return b.orders.BuyOrders.Size()
}

func (b *auction) auctionBookOrders() *BuySellOrders {
	return &b.orders
}

func cancelOrders(bs *BuySellOrders) []*fixmodel.ExecutionReport {
	execs := []*fixmodel.ExecutionReport{}
	for iter := bs.BuyOrders.Iterator(); iter.Next() == true; {
		order := iter.Value().(*orderstate.OrderState)
		bs.BuyOrders.RemoveByID(order.OrderID())
		exec := orderstate.NewOrderCancelledExecutionReport(order)
		execs = append(execs, exec)
	}
	for iter := bs.SellOrders.Iterator(); iter.Next() == true; {
		order := iter.Value().(*orderstate.OrderState)
		bs.SellOrders.RemoveByID(order.OrderID())
		exec := orderstate.NewOrderCancelledExecutionReport(order)
		execs = append(execs, exec)
	}
	return execs
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
