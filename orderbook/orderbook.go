package orderbook

import (
	"fmt"
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"orderbook/instrument"
	"time"
)

type OrderBook interface {
	Instrument() *instrument.Instrument

	NewOrder(order NewOrderSingle) ([]ExecutionReport, error)
	CancelOrder(order OrderCancelRequest) (ExecutionReport, error)
	Tick(time time.Time) ([]ExecutionReport, error)

	OpenTrading() ([]ExecutionReport, error)
	CloseTrading() ([]ExecutionReport, error)
	NoTrading() error

	BuySize() int
	SellSize() int
	BuyOrders() []OrderState
	SellOrders() []OrderState

	orderBookOrders() *buySellOrders

	OrderBookAuction
}

func MakeOrderBook(instrument instrument.Instrument, orderBookEvent OrderBookEventType, clock clock.Clock) OrderBook {
	b := orderbook{instrument: &instrument}
	b.obOrders.buyOrders = NewOrderList(HighToLow)
	b.obOrders.sellOrders = NewOrderList(LowToHigh)
	b.auctionOrders.buyOrders = NewOrderList(HighToLow)
	b.auctionOrders.sellOrders = NewOrderList(LowToHigh)
	b.orderBookState = OrderBookEventTypeAs(orderBookEvent)
	b.clock = clock
	priced, _ := decimal.NewFromString("1.23")
	priced.Add(priced)
	return OrderBook(&b)
}

type buySellOrders struct {
	buyOrders  OrderList
	sellOrders OrderList
}

type orderbook struct {
	instrument     *instrument.Instrument
	obOrders       buySellOrders
	auctionOrders  buySellOrders
	orderBookState OrderBookState
	clock          clock.Clock
}

func (b *orderbook) NewOrder(order NewOrderSingle) ([]ExecutionReport, error) {
	execs := []ExecutionReport{}

	if b.orderBookState == OrderBookStateOrderEntryClosed {
		execs = append(execs, MakeRejectExecutionReport(order))
		return execs, nil
	}

	if b.orderBookState == OrderBookStateTradingClosed && order.TimeInForce() != TimeInForceGoodForAuction {
		execs = append(execs, MakeRejectExecutionReport(order))
		return execs, nil
	}

	if b.orderBookState == OrderBookStateAuctionClosed && order.TimeInForce() == TimeInForceGoodForAuction {
		execs = append(execs, MakeRejectExecutionReport(order))
		return execs, nil
	}

	if order.TimeInForce() == TimeInForceGoodForAuction {
		return addNewOrder(order, &b.auctionOrders, b.clock)
	}

	if b.orderBookState != OrderBookStateTradingOpen {
		if order.OrderType() == OrderTypeMarket || order.TimeInForce() == TimeInForceImmediateOrCancel || order.TimeInForce() == TimeInForceFillOrKill {
			execs = append(execs, MakeRejectExecutionReport(order))
			return execs, nil
		}
		return addNewOrder(order, &b.obOrders, b.clock)
	}

	// reject market orders if there are no limit orders
	if order.OrderType() == OrderTypeMarket {
		if order.isBuy() {
			if b.obOrders.sellOrders.Size() == 0 {
				execs = append(execs, MakeRejectExecutionReport(order))
				return execs, nil
			}
		} else if order.Side() == SideSell {
			if b.obOrders.buyOrders.Size() == 0 {
				execs = append(execs, MakeRejectExecutionReport(order))
				return execs, nil
			}
		}
	}

	execs, err := matchOrderOnBook(order, &b.obOrders, b.clock)

	if order.TimeInForce() == TimeInForceImmediateOrCancel {
		if orderstate := b.obOrders.buyOrders.FindByID(order.OrderID()); orderstate != nil {
			exec := MakeOrderCancelledExecutionReport(orderstate, ExecTypeCanceled)
			execs = append(execs, exec)
			b.obOrders.buyOrders.RemoveByID(orderstate.OrderID())
		}
		if orderstate := b.obOrders.sellOrders.FindByID(order.OrderID()); orderstate != nil {
			exec := MakeOrderCancelledExecutionReport(orderstate, ExecTypeCanceled)
			execs = append(execs, exec)
			b.obOrders.sellOrders.RemoveByID(orderstate.OrderID())
		}
	}
	return execs, err
}

func cancelOrderByFn(ol OrderList, time time.Time, fn func(order OrderState, t time.Time) bool) []ExecutionReport {
	execs := []ExecutionReport{}
	orders := []OrderState{}

	for iter := ol.iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		if fn(order, time) {
			orders = append(orders, order)
			exec := MakeOrderCancelledExecutionReport(order, ExecTypeCanceled)
			execs = append(execs, exec)
		}
	}

	for _, v := range orders {
		ol.RemoveByID(v.OrderID())
	}
	return execs
}

func (b *orderbook) Tick(tm time.Time) ([]ExecutionReport, error) {
	fn := func(order OrderState, t time.Time) bool {
		fmt.Printf("tif %v now %v expire %v\n", order.TimeInForce(), t, order.ExpireOn())
		return order.TimeInForce() == TimeInForceGoodForTime && !t.Before(order.ExpireOn())
	}
	execs := cancelOrderByFn(b.obOrders.buyOrders, tm, fn)
	execs = append(execs, cancelOrderByFn(b.obOrders.sellOrders, tm, fn)...)
	return execs, nil
}

func (b *orderbook) OpenTrading() ([]ExecutionReport, error) {
	var err error
	ex := []ExecutionReport{}
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeOpenTrading)
	if err == nil {
		ex = matchOrder(&b.obOrders)
		//ex = matchOrderReverseSellBuy(&b.obOrders)
	}
	return ex, err
}

func (b *orderbook) CloseTrading() ([]ExecutionReport, error) {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseTrading)
	if err == nil {
		return cancelDayOrders(&b.obOrders, b.clock), nil
	}
	return nil, err
}

func (b *orderbook) orderBookOrders() *buySellOrders {
	return &b.obOrders
}

func cancelDayOrders(bs *buySellOrders, clock clock.Clock) []ExecutionReport {
	fn := func(order OrderState, t time.Time) bool {
		return order.TimeInForce() == TimeInForceDay
	}
	execs := cancelOrderByFn(bs.buyOrders, clock.Now(), fn)
	execs = append(execs, cancelOrderByFn(bs.sellOrders, clock.Now(), fn)...)
	return execs
}

func (b *orderbook) NoTrading() error {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseOrderEntry)
	return err
}

func cancelOrders(bs *buySellOrders) []ExecutionReport {
	execs := []ExecutionReport{}
	for iter := bs.buyOrders.iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		bs.buyOrders.RemoveByID(order.OrderID())
		exec := MakeOrderCancelledExecutionReport(order, ExecTypeCanceled)
		execs = append(execs, exec)
	}
	for iter := bs.sellOrders.iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		bs.sellOrders.RemoveByID(order.OrderID())
		exec := MakeOrderCancelledExecutionReport(order, ExecTypeCanceled)
		execs = append(execs, exec)
	}
	return execs
}

func addNewOrder(order NewOrderSingle, bs *buySellOrders, clock clock.Clock) ([]ExecutionReport, error) {
	execs := []ExecutionReport{}
	var err error
	neworder := NewOrder(order, newID(uuid.NewUUID()), clock.Now())
	execs = append(execs, MakeNewOrderAckExecutionReport(neworder))
	if order.Side() == SideBuy {
		err = bs.buyOrders.Add(neworder)
	} else {
		err = bs.sellOrders.Add(neworder)
	}
	return execs, err
}

func matchOrderOnBook(order NewOrderSingle, bs *buySellOrders, clock clock.Clock) ([]ExecutionReport, error) {
	var err error
	execs := []ExecutionReport{}

	neworder := NewOrder(order, newID(uuid.NewUUID()), clock.Now())
	execs = append(execs, MakeNewOrderAckExecutionReport(neworder))
	//filledBookSellOrders := []OrderState{}
	//filledBookBuyOrders := []OrderState{}

	//fmt.Printf("buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
	matchBuy := func(neworder OrderState, bookorder OrderState) bool {
		return (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() >= bookorder.Price()
	}
	matchSell := func(neworder OrderState, bookorder OrderState) bool {
		return (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() <= bookorder.Price()
	}

	if neworder.TimeInForce() == TimeInForceFillOrKill {
		// check available volume for neworder and if enough carry on otherwise cancel
		if order.isBuy() {
			if !canMatchAnOrder(neworder, bs.sellOrders, bs.buyOrders, execs, matchBuy) {
				execs = append(execs, MakeRejectExecutionReportFromOrderState(neworder))
				return execs, nil
			}
		} else if !order.isBuy() {
			if !canMatchAnOrder(neworder, bs.buyOrders, bs.sellOrders, execs, matchSell) {
				execs = append(execs, MakeRejectExecutionReportFromOrderState(neworder))
				return execs, nil
			}
		}
	}

	if order.isBuy() {
		execs, err = matchAnOrder(neworder, bs.sellOrders, bs.buyOrders, execs, matchBuy)
	} else if !order.isBuy() {
		execs, err = matchAnOrder(neworder, bs.buyOrders, bs.sellOrders, execs, matchSell)
	}

	return execs, err
}

func matchAnOrder(neworder OrderState, bookToMatch OrderList, bookToAdd OrderList, execs []ExecutionReport, match func(neworder OrderState, bookorder OrderState) bool) ([]ExecutionReport, error) {
	filledBookOrders := []OrderState{}
	var err error
	for iter := bookToMatch.iterator(); iter.Next() == true; {
		bookorder := iter.Value().(OrderState)
		if neworder.ClientID() == bookorder.ClientID() {
			execs = append(execs, MakeRejectExecutionReportFromOrderState(neworder))
			return execs, nil
		}
		if match(neworder, bookorder) {
			toFill := min(bookorder.LeavesQty(), neworder.LeavesQty())
			price := bookorder.Price()
			if toFill > 0 {
				neworder.fill(toFill)
				execs = append(execs, MakeFillExecutionReport(neworder, price, toFill))
				if bookorder.fill(toFill) {
					filledBookOrders = append(filledBookOrders, bookorder)
				}
				execs = append(execs, MakeFillExecutionReport(bookorder, price, toFill))
			} else {
				break
			}
		}
		//fmt.Printf("After loop buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
	}
	if neworder.OrderType() == OrderTypeLimit && neworder.LeavesQty() > 0 {
		err = bookToAdd.Add(neworder)
	}
	for _, v := range filledBookOrders {
		bookToMatch.RemoveByID(v.OrderID())
	}
	return execs, err
}

func canMatchAnOrder(neworder OrderState, bookToMatch OrderList, bookToAdd OrderList, execs []ExecutionReport, match func(neworder OrderState, bookorder OrderState) bool) bool {
	var toFill int64
	for iter := bookToMatch.iterator(); iter.Next() == true; {
		bookorder := iter.Value().(OrderState)
		if match(neworder, bookorder) {
			toFill += min(bookorder.LeavesQty(), neworder.LeavesQty())
			if toFill >= neworder.LeavesQty() {
				return true
			}
		}
	}
	return false
}

func matchOrder(bs *buySellOrders) []ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", bs.sellOrders.Size(), bs.buyOrders.Size())
	execs := []ExecutionReport{}
	filledBookSellOrders := []OrderState{}
	filledBookBuyOrders := []OrderState{}

	for buyiter := bs.buyOrders.iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(OrderState)
		//fmt.Printf("buy order %s %f %d\n", SideToString(buyorder.Side()), buyorder.Price(), buyorder.LeavesQty())
		for selliter := bs.sellOrders.iterator(); selliter.Next() == true; {
			sellorder := selliter.Value().(OrderState)
			//fmt.Printf("buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
			if (buyorder.OrderType() == OrderTypeMarket || sellorder.OrderType() == OrderTypeMarket) || buyorder.Price() >= sellorder.Price() {
				toFill := min(sellorder.LeavesQty(), buyorder.LeavesQty())
				var price float64
				if buyorder.OrderType() == OrderTypeMarket {
					price = sellorder.Price()
				} else if sellorder.OrderType() == OrderTypeMarket {
					price = buyorder.Price()
				} else if sellorder.TransactTime().Before(buyorder.TransactTime()) {
					price = sellorder.Price()
				} else {
					price = buyorder.Price()
				}
				if toFill > 0 {
					if buyorder.fill(toFill) {
						filledBookBuyOrders = append(filledBookBuyOrders, buyorder)
					}
					execs = append(execs, MakeFillExecutionReport(buyorder, price, toFill))
					if sellorder.fill(toFill) {
						filledBookSellOrders = append(filledBookSellOrders, sellorder)
					}
					execs = append(execs, MakeFillExecutionReport(sellorder, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
		}
	}
	// remove filled orders
	for _, v := range filledBookBuyOrders {
		bs.buyOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBookSellOrders {
		bs.sellOrders.RemoveByID(v.OrderID())
	}

	return execs
}

func matchOrderReverseSellBuy(bs *buySellOrders) []ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", bs.sellOrders.Size(), bs.buyOrders.Size())
	execs := []ExecutionReport{}
	filledBookbuy2orders := []OrderState{}
	filledBooksell1orders := []OrderState{}

	for iter1 := bs.sellOrders.iterator(); iter1.Next() == true; {
		sell1order := iter1.Value().(OrderState)
		//fmt.Printf("buy order %s %f %d\n", SideToString(sell1order.Side()), sell1order.Price(), sell1order.LeavesQty())
		for iter2 := bs.buyOrders.iterator(); iter2.Next() == true; {
			buy2order := iter2.Value().(OrderState)
			//fmt.Printf("buy \nbuy2order %v \nbuyorder %v\n", buy2order, buyorder)
			if (sell1order.OrderType() == OrderTypeMarket || buy2order.OrderType() == OrderTypeMarket) || sell1order.Price() <= buy2order.Price() {
				toFill := min(buy2order.LeavesQty(), sell1order.LeavesQty())
				var price float64
				if sell1order.OrderType() == OrderTypeMarket {
					price = buy2order.Price()
				} else if buy2order.OrderType() == OrderTypeMarket {
					price = sell1order.Price()
				} else if sell1order.TransactTime().Before(buy2order.TransactTime()) {
					price = sell1order.Price()
				} else {
					price = buy2order.Price()
				}
				if toFill > 0 {
					if sell1order.fill(toFill) {
						filledBooksell1orders = append(filledBooksell1orders, sell1order)
					}
					execs = append(execs, MakeFillExecutionReport(sell1order, price, toFill))
					if buy2order.fill(toFill) {
						filledBookbuy2orders = append(filledBookbuy2orders, buy2order)
					}
					execs = append(execs, MakeFillExecutionReport(buy2order, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nbuy2order %v \nsell1order %v\n", buy2order, sell1order)
		}
	}
	// remove filled orders
	for _, v := range filledBooksell1orders {
		bs.sellOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBookbuy2orders {
		bs.buyOrders.RemoveByID(v.OrderID())
	}

	return execs
}

func (b *orderbook) CancelOrder(order OrderCancelRequest) (ExecutionReport, error) {
	var exec, err = cancelOrder(order, &b.obOrders)
	// try auction orders if not found on orderbook
	if exec == nil {
		exec, err = cancelOrder(order, &b.auctionOrders)
	}
	return exec, err
}

func cancelOrder(order OrderCancelRequest, bs *buySellOrders) (ExecutionReport, error) {
	var ord OrderState
	if order.Side() == SideBuy {
		ord = bs.buyOrders.FindByClOrdID(order.OrigClOrdID())
		if ord != nil {
			bs.buyOrders.RemoveByID(ord.OrderID())
			return MakeCancelOrderExecutionReport(ord, order), nil
		}
	} else {
		ord = bs.sellOrders.FindByClOrdID(order.OrigClOrdID())
		if ord != nil {
			bs.sellOrders.RemoveByID(ord.OrderID())
			return MakeCancelOrderExecutionReport(ord, order), nil
		}
	}
	return nil, nil
}

func (b *orderbook) Instrument() *instrument.Instrument {
	return b.instrument
}

func (b *orderbook) BuySize() int {
	return b.obOrders.buyOrders.Size()
}

func (b *orderbook) SellSize() int {
	return b.obOrders.sellOrders.Size()
}

func (b *orderbook) BuyOrders() []OrderState {
	return b.obOrders.buyOrders.Orders()
}

func (b *orderbook) SellOrders() []OrderState {
	return b.obOrders.sellOrders.Orders()
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
