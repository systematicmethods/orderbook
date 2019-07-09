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

	NewAuctionOrder(order NewOrderSingle) ([]ExecutionReport, error)
	Auction() ([]ExecutionReport, error)

	matchOrder() []ExecutionReport
	addNewOrder(order NewOrderSingle) ([]ExecutionReport, error)
	matchOrderOnBook(order NewOrderSingle) ([]ExecutionReport, error)
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
	//execs, _ := b.addNewOrder(order)
	execs, _ := b.matchOrderOnBook(order)
	//execs = append(execs, b.matchOrder()...)
	return execs, nil
}

func (b *orderbook) NewAuctionOrder(order NewOrderSingle) ([]ExecutionReport, error) {
	return b.addNewOrder(order)
}

func (b *orderbook) Auction() ([]ExecutionReport, error) {
	return b.matchOrder(), nil
}

func (b *orderbook) addNewOrder(order NewOrderSingle) ([]ExecutionReport, error) {
	execs := []ExecutionReport{}
	if order.OrderType() == OrderTypeMarket {
		if order.Side() == SideBuy {
			if b.sellOrders.Size() == 0 {
				execs = append(execs, MakeRejectExecutionReport(order))
				return execs, nil
			}
		} else if order.Side() == SideSell {
			if b.buyOrders.Size() == 0 {
				execs = append(execs, MakeRejectExecutionReport(order))
				return execs, nil
			}
		}
	}
	if order.OrderID() != "" {
		order := NewOrder(order, newID(uuid.NewUUID()), time.Now())
		//fmt.Printf("NewOrder added %v added qty %d\n", order, order.OrderQty())
		execs = append(execs, MakeNewOrderAckExecutionReport(order))
		//fmt.Printf("NewOrder execs in OrderBook %v\n", execs)
		if order.Side() == SideBuy {
			b.buyOrders.Add(order)
		} else {
			b.sellOrders.Add(order)
		}
		//fmt.Printf("NewOrder execs after in OrderBook %v", execs)
		return execs, nil
	}
	return nil, nil
}

func (b *orderbook) matchOrderOnBook(order NewOrderSingle) ([]ExecutionReport, error) {
	matchexecs := []ExecutionReport{}

	// reject market orders if there are no limit orders
	if order.OrderType() == OrderTypeMarket {
		if order.isBuy() {
			if b.sellOrders.Size() == 0 {
				matchexecs = append(matchexecs, MakeRejectExecutionReport(order))
				return matchexecs, nil
			}
		} else if order.Side() == SideSell {
			if b.buyOrders.Size() == 0 {
				matchexecs = append(matchexecs, MakeRejectExecutionReport(order))
				return matchexecs, nil
			}
		}
	}

	neworder := NewOrder(order, newID(uuid.NewUUID()), time.Now())
	matchexecs = append(matchexecs, MakeNewOrderAckExecutionReport(neworder))
	filledBookSellOrders := []OrderState{}
	filledBookBuyOrders := []OrderState{}

	if order.isBuy() && b.sellOrders.Size() > 0 {
		for iter := b.sellOrders.iterator(); iter.Next() == true; {
			bookorder := iter.Value().(OrderState)
			if order.ClientID() == bookorder.ClientID() {
				matchexecs = append(matchexecs, MakeRejectExecutionReport(order))
				return matchexecs, nil
			}
			//fmt.Printf("buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
			if (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() >= bookorder.Price() {
				toFill := min(bookorder.LeavesQty(), neworder.LeavesQty())
				price := bookorder.Price()
				if toFill > 0 {
					neworder.fill(toFill)
					matchexecs = append(matchexecs, MakeFillExecutionReport(neworder, price, toFill))
					if bookorder.fill(toFill) {
						filledBookSellOrders = append(filledBookSellOrders, bookorder)
					}
					matchexecs = append(matchexecs, MakeFillExecutionReport(bookorder, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
		}
	} else if !order.isBuy() && b.buyOrders.Size() > 0 {
		for iter := b.buyOrders.iterator(); iter.Next() == true; {
			bookorder := iter.Value().(OrderState)
			if order.ClientID() == bookorder.ClientID() {
				matchexecs = append(matchexecs, MakeRejectExecutionReport(order))
				return matchexecs, nil
			}
			//fmt.Printf("buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
			if (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() <= bookorder.Price() {
				toFill := min(bookorder.LeavesQty(), neworder.LeavesQty())
				price := bookorder.Price()
				if toFill > 0 {
					neworder.fill(toFill)
					matchexecs = append(matchexecs, MakeFillExecutionReport(neworder, price, toFill))
					if bookorder.fill(toFill) {
						filledBookBuyOrders = append(filledBookBuyOrders, bookorder)
					}
					matchexecs = append(matchexecs, MakeFillExecutionReport(bookorder, price, toFill))
				} else {
					break
				}
			}
			//fmt.Printf("After loop buy \nbookorder %v \nneworder %v\n", bookorder, neworder)
		}
	}
	if neworder.OrderType() == OrderTypeLimit && neworder.LeavesQty() > 0 {
		if order.isBuy() {
			b.buyOrders.Add(neworder)
		} else {
			b.sellOrders.Add(neworder)
		}
	}

	// remove filled orders
	for _, v := range filledBookBuyOrders {
		b.buyOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBookSellOrders {
		b.sellOrders.RemoveByID(v.OrderID())
	}
	return matchexecs, nil
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
							b.buyOrders.RemoveByID(buyorder.OrderID())
						}
						matchexecs = append(matchexecs, MakeFillExecutionReport(buyorder, price, toFill))
						if sellorder.fill(toFill) {
							b.sellOrders.RemoveByID(sellorder.OrderID())
						}
						matchexecs = append(matchexecs, MakeFillExecutionReport(sellorder, price, toFill))
					} else {
						break
					}
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
