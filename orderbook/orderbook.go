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

	OpenTrading() ([]ExecutionReport, error)
	CloseTrading() error
	NoTrading() error
	OpenAuction() error
	CloseAuction() ([]ExecutionReport, error)

	BuySize() int
	SellSize() int
	BuyOrders() []OrderState
	SellOrders() []OrderState

	BuyAuctionSize() int
	SellAuctionSize() int
	BuyAuctionOrders() []OrderState
	SellAuctionOrders() []OrderState
}

func MakeOrderBook(instrument instrument.Instrument, orderBookEvent OrderBookEventType) OrderBook {
	b := orderbook{instrument: &instrument}
	b.buyOrders = NewOrderList(HighToLow)
	b.sellOrders = NewOrderList(HighToLow)
	b.obOrders.buyOrders = NewOrderList(HighToLow)
	b.obOrders.sellOrders = NewOrderList(HighToLow)
	b.auctionOrders.buyOrders = NewOrderList(HighToLow)
	b.auctionOrders.sellOrders = NewOrderList(HighToLow)
	b.orderBookState = OrderBookEventTypeAs(orderBookEvent)
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
	buyOrders      OrderList
	sellOrders     OrderList
	orderBookState OrderBookState
}

func (b *orderbook) NewOrder(order NewOrderSingle) ([]ExecutionReport, error) {
	execs := []ExecutionReport{}
	// reject all orders when trading not open
	if b.orderBookState != OrderBookStateTradingOpen {
		if b.orderBookState == OrderBookStateAuctionOpen {
			return addNewOrder(order, &b.auctionOrders)
		}
		execs = append(execs, MakeRejectExecutionReport(order))
		return execs, nil
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

	execs, err := matchOrderOnBook(order, &b.obOrders)
	return execs, err
}

func (b *orderbook) OpenTrading() ([]ExecutionReport, error) {
	var err error
	ex := []ExecutionReport{}
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeOpenTrading)
	if err == nil {
		ex = matchOrder(&b.obOrders)
	}
	return ex, err
}

func (b *orderbook) CloseTrading() error {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseTrading)
	return err
}

func (b *orderbook) NoTrading() error {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeNoTrading)
	return err
}

func (b *orderbook) OpenAuction() error {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeOpenAuction)
	return err
}

func (b *orderbook) CloseAuction() ([]ExecutionReport, error) {
	var err error
	ex := []ExecutionReport{}
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseAuction)
	if err == nil {
		ex = matchOrdersOnBook(&b.auctionOrders)
		// TODO: Cancel remaining action orders in the auction
	}
	return ex, err
}

func addNewOrder(order NewOrderSingle, bs *buySellOrders) ([]ExecutionReport, error) {
	execs := []ExecutionReport{}
	var err error
	neworder := NewOrder(order, newID(uuid.NewUUID()), time.Now())
	execs = append(execs, MakeNewOrderAckExecutionReport(neworder))
	if order.Side() == SideBuy {
		err = bs.buyOrders.Add(neworder)
	} else {
		err = bs.sellOrders.Add(neworder)
	}
	return execs, err
}

func matchOrderOnBook(order NewOrderSingle, bs *buySellOrders) ([]ExecutionReport, error) {
	var err error
	execs := []ExecutionReport{}

	neworder := NewOrder(order, newID(uuid.NewUUID()), time.Now())
	execs = append(execs, MakeNewOrderAckExecutionReport(neworder))
	filledBookSellOrders := []OrderState{}
	filledBookBuyOrders := []OrderState{}

	if order.isBuy() && bs.sellOrders.Size() > 0 {
		for iter := bs.sellOrders.iterator(); iter.Next() == true; {
			bookorder := iter.Value().(OrderState)
			if order.ClientID() == bookorder.ClientID() {
				execs = append(execs, MakeRejectExecutionReport(order))
				return execs, nil
			}
			//fmt.Printf("buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
			if (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() >= bookorder.Price() {
				toFill := min(bookorder.LeavesQty(), neworder.LeavesQty())
				price := bookorder.Price()
				if toFill > 0 {
					neworder.fill(toFill)
					execs = append(execs, MakeFillExecutionReport(neworder, price, toFill))
					if bookorder.fill(toFill) {
						filledBookSellOrders = append(filledBookSellOrders, bookorder)
					}
					execs = append(execs, MakeFillExecutionReport(bookorder, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
		}
	} else if !order.isBuy() && bs.buyOrders.Size() > 0 {
		for iter := bs.buyOrders.iterator(); iter.Next() == true; {
			bookorder := iter.Value().(OrderState)
			if order.ClientID() == bookorder.ClientID() {
				execs = append(execs, MakeRejectExecutionReport(order))
				return execs, nil
			}
			//fmt.Printf("buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
			if (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() <= bookorder.Price() {
				toFill := min(bookorder.LeavesQty(), neworder.LeavesQty())
				price := bookorder.Price()
				if toFill > 0 {
					neworder.fill(toFill)
					execs = append(execs, MakeFillExecutionReport(neworder, price, toFill))
					if bookorder.fill(toFill) {
						filledBookBuyOrders = append(filledBookBuyOrders, bookorder)
					}
					execs = append(execs, MakeFillExecutionReport(bookorder, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
		}
	}
	if neworder.OrderType() == OrderTypeLimit && neworder.LeavesQty() > 0 {
		if order.isBuy() {
			err = bs.buyOrders.Add(neworder)
		} else {
			err = bs.sellOrders.Add(neworder)
		}
	}

	// remove filled orders
	for _, v := range filledBookBuyOrders {
		bs.buyOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBookSellOrders {
		bs.sellOrders.RemoveByID(v.OrderID())
	}
	return execs, err
}

func matchOrder(bs *buySellOrders) []ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", bs.sellOrders.Size(), bs.buyOrders.Size())
	execs := []ExecutionReport{}
	for buyiter := bs.buyOrders.iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(OrderState)
		if buyorder.Side() == SideBuy && bs.sellOrders.Size() > 0 {
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
					} else {
						price = buyorder.Price()
					}
					if toFill > 0 {
						if buyorder.fill(toFill) {
							bs.buyOrders.RemoveByID(buyorder.OrderID())
						}
						execs = append(execs, MakeFillExecutionReport(buyorder, price, toFill))
						if sellorder.fill(toFill) {
							bs.sellOrders.RemoveByID(sellorder.OrderID())
						}
						execs = append(execs, MakeFillExecutionReport(sellorder, price, toFill))
					} else {
						break
					}
				}
				//fmt.Printf("After loop buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
			}
		}
	}
	return execs
}

func matchOrdersOnBook(bs *buySellOrders) []ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", b.sellOrders.Size(), b.buyOrders.Size())
	execs := []ExecutionReport{}
	for buyiter := bs.buyOrders.iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(OrderState)
		if buyorder.Side() == SideBuy && bs.sellOrders.Size() > 0 {
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
					} else {
						price = buyorder.Price()
					}
					if toFill > 0 {
						if buyorder.fill(toFill) {
							bs.buyOrders.RemoveByID(buyorder.OrderID())
						}
						execs = append(execs, MakeFillExecutionReport(buyorder, price, toFill))
						if sellorder.fill(toFill) {
							bs.sellOrders.RemoveByID(sellorder.OrderID())
						}
						execs = append(execs, MakeFillExecutionReport(sellorder, price, toFill))
					} else {
						break
					}
				}
				//fmt.Printf("After loop buy \nsellorder %v \nbuyorder %v\n", sellorder, buyorder)
			}
		}
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

func (b *orderbook) BuyAuctionSize() int {
	return b.auctionOrders.buyOrders.Size()
}

func (b *orderbook) SellAuctionSize() int {
	return b.auctionOrders.sellOrders.Size()
}

func (b *orderbook) BuyAuctionOrders() []OrderState {
	return b.auctionOrders.buyOrders.Orders()
}

func (b *orderbook) SellAuctionOrders() []OrderState {
	return b.auctionOrders.sellOrders.Orders()
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
