package orderbook

import (
	"fmt"
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"orderbook/etype"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/obmath"
	"orderbook/orderstate"
	"orderbook/tradingevent"
	"time"
)

const (
	rejected = etype.Error("rejected")
)

type OrderBook interface {
	Instrument() *instrument.Instrument

	NewOrder(order *fixmodel.NewOrderSingle) ([]*fixmodel.ExecutionReport, error)
	CancelOrder(order *fixmodel.OrderCancelRequest) (*fixmodel.ExecutionReport, error)
	Tick(time time.Time) ([]*fixmodel.ExecutionReport, error)

	OpenTrading() ([]*fixmodel.ExecutionReport, error)
	CloseTrading() ([]*fixmodel.ExecutionReport, error)
	NoTrading() error

	BuyOrders() []*orderstate.OrderState
	SellOrders() []*orderstate.OrderState

	BuySize() int
	SellSize() int

	BuyTop() *MarketPrice
	SellTop() *MarketPrice

	State() tradingevent.OrderBookState
}

type MarketPrice struct {
	Price    float64
	Quantity int64
}

type BuySellOrders struct {
	BuyOrders  *orderstate.Orderlist
	SellOrders *orderstate.Orderlist
}

type orderbook struct {
	instrument     *instrument.Instrument
	orders         BuySellOrders
	orderBookState tradingevent.OrderBookState
	clock          clock.Clock
}

func (b *orderbook) State() tradingevent.OrderBookState {
	return b.orderBookState
}

func (b *orderbook) NewOrder(order *fixmodel.NewOrderSingle) ([]*fixmodel.ExecutionReport, error) {
	execs := []*fixmodel.ExecutionReport{}

	if order.TimeInForce() == fixmodel.TimeInForceGoodForAuction {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonUnsupportedOrderCharacteristic, "Auction orders not allowed in limit order books"))
		return execs, rejected
	}

	if b.State() == tradingevent.OrderBookStateAuctionOpen {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, "Auction open rejecting non auction orders"))
		return execs, rejected
	}

	if b.State() == tradingevent.OrderBookStateOrderEntryClosed {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, "Exchange closed"))
		return execs, rejected
	}

	if b.State() == tradingevent.OrderBookStateTradingClosed {
		execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonExchangeClosed, "Trading closed"))
		return execs, rejected
	}

	//if b.orderBookState == tradingevent.OrderBookStateAuctionClosed && order.TimeInForce() == fixmodel.TimeInForceGoodForAuction {
	//	execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonDuplicateOrder, ""))
	//	return execs, nil
	//}

	if b.State() != tradingevent.OrderBookStateTradingOpen {
		if order.OrderType() == fixmodel.OrderTypeMarket || order.TimeInForce() == fixmodel.TimeInForceImmediateOrCancel || order.TimeInForce() == fixmodel.TimeInForceFillOrKill {
			execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonOther, "trading not open"))
			return execs, rejected
		}
		return addNewOrder(order, &b.orders, b.clock)
	}

	// reject market orders if there are no limit orders
	if order.OrderType() == fixmodel.OrderTypeMarket {
		if order.Side() == fixmodel.SideBuy {
			if b.orders.SellOrders.Size() == 0 {
				execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonOther, "empty order book"))
				return execs, rejected
			}
		} else if order.Side() == fixmodel.SideSell {
			if b.orders.BuyOrders.Size() == 0 {
				execs = append(execs, fixmodel.NewRejectExecutionReport(order, fixmodel.OrdRejReasonOther, "mpty order book"))
				return execs, rejected
			}
		}
	}
	//// add stop orders to the stoplist
	//if order.OrderType() == OrderTypeStop {
	//	return addNewOrder(order, &b.stopOrders, b.clock)
	//}

	execs, err := matchOrderOnBook(order, &b.orders, b.clock)

	if order.TimeInForce() == fixmodel.TimeInForceImmediateOrCancel {
		if order := b.orders.BuyOrders.FindByID(order.OrderID()); order != nil {
			exec := orderstate.NewOrderCancelledExecutionReport(order)
			execs = append(execs, exec)
			b.orders.BuyOrders.RemoveByID(order.OrderID())
		}
		if order := b.orders.SellOrders.FindByID(order.OrderID()); order != nil {
			exec := orderstate.NewOrderCancelledExecutionReport(order)
			execs = append(execs, exec)
			b.orders.SellOrders.RemoveByID(order.OrderID())
		}
	}
	return execs, err
}

func cancelOrderByFn(ol *orderstate.Orderlist, time time.Time, fn func(order *orderstate.OrderState, t time.Time) bool) []*fixmodel.ExecutionReport {
	execs := []*fixmodel.ExecutionReport{}
	orders := []*orderstate.OrderState{}

	for iter := ol.Iterator(); iter.Next() == true; {
		order := iter.Value().(*orderstate.OrderState)
		if fn(order, time) {
			orders = append(orders, order)
			exec := orderstate.NewOrderCancelledExecutionReport(order)
			execs = append(execs, exec)
		}
	}

	for _, v := range orders {
		ol.RemoveByID(v.OrderID())
	}
	return execs
}

func (b *orderbook) Tick(tm time.Time) ([]*fixmodel.ExecutionReport, error) {
	fn := func(order *orderstate.OrderState, t time.Time) bool {
		fmt.Printf("tif %v now %v expire %v\n", order.TimeInForce(), t, order.ExpireOn())
		return order.TimeInForce() == fixmodel.TimeInForceGoodForTime && !t.Before(order.ExpireOn())
	}
	execs := cancelOrderByFn(b.orders.BuyOrders, tm, fn)
	execs = append(execs, cancelOrderByFn(b.orders.SellOrders, tm, fn)...)
	return execs, nil
}

func (b *orderbook) OpenTrading() ([]*fixmodel.ExecutionReport, error) {
	var err error
	ex := []*fixmodel.ExecutionReport{}
	b.orderBookState, err = tradingevent.OrderBookStateChange(b.State(), tradingevent.OrderBookEventTypeOpenTrading)
	if err == nil {
		ex = matchOrder(&b.orders)
		//ex = matchOrderReverseSellBuy(&b.orders)
	}
	return ex, err
}

func (b *orderbook) CloseTrading() ([]*fixmodel.ExecutionReport, error) {
	var err error
	b.orderBookState, err = tradingevent.OrderBookStateChange(b.State(), tradingevent.OrderBookEventTypeCloseTrading)
	if err == nil {
		return cancelDayOrders(&b.orders, b.clock), nil
	}
	return nil, err
}

func cancelDayOrders(bs *BuySellOrders, clock clock.Clock) []*fixmodel.ExecutionReport {
	fn := func(order *orderstate.OrderState, t time.Time) bool {
		return order.TimeInForce() == fixmodel.TimeInForceDay
	}
	execs := cancelOrderByFn(bs.BuyOrders, clock.Now(), fn)
	execs = append(execs, cancelOrderByFn(bs.SellOrders, clock.Now(), fn)...)
	return execs
}

func (b *orderbook) NoTrading() error {
	var err error
	b.orderBookState, err = tradingevent.OrderBookStateChange(b.State(), tradingevent.OrderBookEventTypeCloseOrderEntry)
	return err
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

func matchOrderOnBook(order *fixmodel.NewOrderSingle, bs *BuySellOrders, clock clock.Clock) ([]*fixmodel.ExecutionReport, error) {
	var err error
	execs := []*fixmodel.ExecutionReport{}

	neworder := orderstate.NewOrder(order, newID(uuid.NewUUID()), clock.Now())
	//fmt.Println("neworder", neworder)
	execs = append(execs, orderstate.NewNewOrderAckExecutionReport(neworder))
	//filledBookSellOrders := []fixmodel.OrderState{}
	//filledBookBuyOrders := []fixmodel.OrderState{}

	//fmt.Printf("buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
	matchBuy := func(neworder *orderstate.OrderState, bookorder *orderstate.OrderState) bool {
		return (neworder.OrderType() == fixmodel.OrderTypeMarket || bookorder.OrderType() == fixmodel.OrderTypeMarket) || neworder.Price() >= bookorder.Price()
	}
	matchSell := func(neworder *orderstate.OrderState, bookorder *orderstate.OrderState) bool {
		return (neworder.OrderType() == fixmodel.OrderTypeMarket || bookorder.OrderType() == fixmodel.OrderTypeMarket) || neworder.Price() <= bookorder.Price()
	}

	if neworder.TimeInForce() == fixmodel.TimeInForceFillOrKill {
		// check available volume for neworder and if enough carry on otherwise cancel
		if order.Side() == fixmodel.SideBuy {
			if !canMatchAnOrder(neworder, bs.SellOrders, bs.BuyOrders, execs, matchBuy) {
				execs = append(execs, orderstate.NewRejectExecutionRepor(neworder))
				return execs, nil
			}
		} else if order.Side() == fixmodel.SideSell {
			if !canMatchAnOrder(neworder, bs.BuyOrders, bs.SellOrders, execs, matchSell) {
				execs = append(execs, orderstate.NewRejectExecutionRepor(neworder))
				return execs, nil
			}
		}
	}

	if order.Side() == fixmodel.SideBuy {
		execs, err = matchAnOrder(neworder, bs.SellOrders, bs.BuyOrders, execs, matchBuy)
	} else if order.Side() == fixmodel.SideSell {
		execs, err = matchAnOrder(neworder, bs.BuyOrders, bs.SellOrders, execs, matchSell)
	}

	return execs, err
}

func matchAnOrder(neworder *orderstate.OrderState, bookToMatch *orderstate.Orderlist, bookToAdd *orderstate.Orderlist, execs []*fixmodel.ExecutionReport,
	match func(neworder *orderstate.OrderState, bookorder *orderstate.OrderState) bool) ([]*fixmodel.ExecutionReport, error) {
	filledBookOrders := []*orderstate.OrderState{}
	var err error
	for iter := bookToMatch.Iterator(); iter.Next() == true; {
		bookorder := iter.Value().(*orderstate.OrderState)
		if neworder.ClientID() == bookorder.ClientID() {
			execs = append(execs, orderstate.NewRejectExecutionRepor(neworder))
			return execs, nil
		}
		if match(neworder, bookorder) {
			toFill := obmath.Min(bookorder.LeavesQty(), neworder.LeavesQty())
			price := bookorder.Price()
			if toFill > 0 {
				neworder.Fill(toFill)
				execs = append(execs, orderstate.NewFillExecutionReport(neworder, price, toFill))
				if bookorder.Fill(toFill) {
					filledBookOrders = append(filledBookOrders, bookorder)
				}
				execs = append(execs, orderstate.NewFillExecutionReport(bookorder, price, toFill))
			} else {
				break
			}
		}
		//fmt.Printf("After loop buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
	}
	if neworder.OrderType() == fixmodel.OrderTypeLimit && neworder.LeavesQty() > 0 {
		err = bookToAdd.Add(neworder)
	}
	for _, v := range filledBookOrders {
		bookToMatch.RemoveByID(v.OrderID())
	}
	return execs, err
}

func canMatchAnOrder(neworder *orderstate.OrderState, bookToMatch *orderstate.Orderlist, bookToAdd *orderstate.Orderlist, execs []*fixmodel.ExecutionReport, match func(neworder *orderstate.OrderState, bookorder *orderstate.OrderState) bool) bool {
	var toFill int64
	for iter := bookToMatch.Iterator(); iter.Next() == true; {
		bookorder := iter.Value().(*orderstate.OrderState)
		if match(neworder, bookorder) {
			toFill += obmath.Min(bookorder.LeavesQty(), neworder.LeavesQty())
			if toFill >= neworder.LeavesQty() {
				return true
			}
		}
	}
	return false
}

func matchOrder(bs *BuySellOrders) []*fixmodel.ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", bs.SellOrders.Size(), bs.BuyOrders.Size())
	execs := []*fixmodel.ExecutionReport{}
	filledBookSellOrders := []*orderstate.OrderState{}
	filledBookBuyOrders := []*orderstate.OrderState{}

	for buyiter := bs.BuyOrders.Iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(*orderstate.OrderState)
		//fmt.Printf("buy order %s %f %d\n", SideToString(buyorder.Side()), buyorder.Price(), buyorder.LeavesQty())
		for selliter := bs.SellOrders.Iterator(); selliter.Next() == true; {
			sellorder := selliter.Value().(*orderstate.OrderState)
			//fmt.Printf("buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
			if (buyorder.OrderType() == fixmodel.OrderTypeMarket || sellorder.OrderType() == fixmodel.OrderTypeMarket) || buyorder.Price() >= sellorder.Price() {
				toFill := obmath.Min(sellorder.LeavesQty(), buyorder.LeavesQty())
				var price float64
				if buyorder.OrderType() == fixmodel.OrderTypeMarket {
					price = sellorder.Price()
				} else if sellorder.OrderType() == fixmodel.OrderTypeMarket {
					price = buyorder.Price()
				} else if sellorder.TransactTime().Before(buyorder.TransactTime()) {
					price = sellorder.Price()
				} else {
					price = buyorder.Price()
				}
				if toFill > 0 {
					if buyorder.Fill(toFill) {
						filledBookBuyOrders = append(filledBookBuyOrders, buyorder)
					}
					execs = append(execs, orderstate.NewFillExecutionReport(buyorder, price, toFill))
					if sellorder.Fill(toFill) {
						filledBookSellOrders = append(filledBookSellOrders, sellorder)
					}
					execs = append(execs, orderstate.NewFillExecutionReport(sellorder, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
		}
	}
	// remove filled orders
	for _, v := range filledBookBuyOrders {
		bs.BuyOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBookSellOrders {
		bs.SellOrders.RemoveByID(v.OrderID())
	}

	return execs
}

func matchOrderReverseSellBuy(bs *BuySellOrders) []*fixmodel.ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", bs.SellOrders.Size(), bs.BuyOrders.Size())
	execs := []*fixmodel.ExecutionReport{}
	filledBookbuy2orders := []*orderstate.OrderState{}
	filledBooksell1orders := []*orderstate.OrderState{}

	for iter1 := bs.SellOrders.Iterator(); iter1.Next() == true; {
		sell1order := iter1.Value().(*orderstate.OrderState)
		//fmt.Printf("buy order %s %f %d\n", SideToString(sell1order.Side()), sell1order.Price(), sell1order.LeavesQty())
		for iter2 := bs.BuyOrders.Iterator(); iter2.Next() == true; {
			buy2order := iter2.Value().(*orderstate.OrderState)
			//fmt.Printf("buy \nbuy2order %v \nbuyorder %v\n", buy2order, buyorder)
			if (sell1order.OrderType() == fixmodel.OrderTypeMarket || buy2order.OrderType() == fixmodel.OrderTypeMarket) || sell1order.Price() <= buy2order.Price() {
				toFill := obmath.Min(buy2order.LeavesQty(), sell1order.LeavesQty())
				var price float64
				if sell1order.OrderType() == fixmodel.OrderTypeMarket {
					price = buy2order.Price()
				} else if buy2order.OrderType() == fixmodel.OrderTypeMarket {
					price = sell1order.Price()
				} else if sell1order.TransactTime().Before(buy2order.TransactTime()) {
					price = sell1order.Price()
				} else {
					price = buy2order.Price()
				}
				if toFill > 0 {
					if sell1order.Fill(toFill) {
						filledBooksell1orders = append(filledBooksell1orders, sell1order)
					}
					execs = append(execs, orderstate.NewFillExecutionReport(sell1order, price, toFill))
					if buy2order.Fill(toFill) {
						filledBookbuy2orders = append(filledBookbuy2orders, buy2order)
					}
					execs = append(execs, orderstate.NewFillExecutionReport(buy2order, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nbuy2order %v \nsell1order %v\n", buy2order, sell1order)
		}
	}
	// remove filled orders
	for _, v := range filledBooksell1orders {
		bs.SellOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBookbuy2orders {
		bs.BuyOrders.RemoveByID(v.OrderID())
	}

	return execs
}

func (b *orderbook) CancelOrder(order *fixmodel.OrderCancelRequest) (*fixmodel.ExecutionReport, error) {
	var exec, err = cancelOrder(order, &b.orders)
	// try auction orders if not found on orderbook
	//if exec == nil {
	//	exec, err = cancelOrder(order, &b.auctionOrders)
	//}
	return exec, err
}

func cancelOrder(order *fixmodel.OrderCancelRequest, bs *BuySellOrders) (*fixmodel.ExecutionReport, error) {
	if order.Side() == fixmodel.SideBuy {
		//ord = bs.BuyOrders.FindByClOrdID(order.OrigClOrdID())
		ord := bs.BuyOrders.FindFirst(func(ord interface{}) bool {
			return ord.(*orderstate.OrderState).ClOrdID() == order.OrigClOrdID()
		})
		if ord != nil {
			bs.BuyOrders.RemoveByID(ord.OrderID())
			return orderstate.NewCancelOrderExecutionReport(ord, order), nil
		}
	} else {
		//ord = bs.SellOrders.FindByClOrdID(order.OrigClOrdID())
		ord := bs.SellOrders.FindFirst(func(ord interface{}) bool {
			return ord.(*orderstate.OrderState).ClOrdID() == order.OrigClOrdID()
		})
		if ord != nil {
			bs.SellOrders.RemoveByID(ord.OrderID())
			return orderstate.NewCancelOrderExecutionReport(ord, order), nil
		}
	}
	return nil, nil
}

//func findByClOrdID(ord interface{}) bool {
//	return ord.(*fixmodel.OrderState).Price() == 1.1
//}

func (b *orderbook) Instrument() *instrument.Instrument {
	return b.instrument
}

func (b *orderbook) BuySize() int {
	return b.orders.BuyOrders.Size()
}

func (b *orderbook) SellSize() int {
	return b.orders.SellOrders.Size()
}

func (b *orderbook) BuyTop() *MarketPrice {
	return &MarketPrice{b.orders.BuyOrders.Top().Price(), b.orders.BuyOrders.Top().LeavesQty()}
}

func (b *orderbook) SellTop() *MarketPrice {
	return &MarketPrice{b.orders.SellOrders.Top().Price(), b.orders.SellOrders.Top().LeavesQty()}
}

func (b *orderbook) BuyOrders() []*orderstate.OrderState {
	ords := []*orderstate.OrderState{}
	xx := b.orders.BuyOrders.Orders2(ords)
	return xx
}

func (b *orderbook) SellOrders() []*orderstate.OrderState {
	ords := []*orderstate.OrderState{}
	return b.orders.SellOrders.Orders2(ords)
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}
