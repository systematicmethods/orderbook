package orderbook

import "math"

func (b *orderbook) OpenAuction() error {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeOpenAuction)
	return err
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
func (b *orderbook) CloseAuction() ([]ExecutionReport, error) {
	var err error
	execs := []ExecutionReport{}
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseAuction)
	if err == nil {
		var exs = matchAuctionOrdersOnBook(&b.auctionOrders)
		execs = append(execs, exs...)
		exs = cancelOrders(&b.auctionOrders)
		execs = append(execs, exs...)
	}
	return execs, err
}

func (b *orderbook) auctionBookOrders() *buySellOrders {
	return &b.auctionOrders
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

func matchAuctionOrdersOnBookMaxVolume(bs *buySellOrders) []ExecutionReport {
	//fmt.Printf("match order: sell %d buy %d \n", b.sellOrders.Size(), b.buyOrders.Size())
	execs := []ExecutionReport{}

	priceb, volb, errb := minPriceOnBuySide(bs)
	prices, vols, errs := maxPriceOnSellSide(bs)

	minvol := math.Min(float64(volb), float64(vols))
	println("minvol", minvol, priceb, prices, errb, errs)
	return execs
}

func matchAuctionOrdersOnBook(bs *buySellOrders) []ExecutionReport {
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

func volMaxPriceMinPriceMax(bs *buySellOrders) (maxvol int64, pricemin float64, pricemax float64, err error) {
	pricemin, volb, errmin := minPriceOnBuySide(bs)
	pricemax, vols, errmax := maxPriceOnSellSide(bs)
	maxvol = min(volb, vols)
	if errmax != nil {
		err = errmax
	}
	if errmin != nil {
		err = errmin
	}
	return
}

func buyVWAP(orders OrderList, maxvol int64) (vwap float64) {
	for iter := orders.iterator(); iter.Next() == true; {

	}
	return
}
func maxPriceOnSellSide(bs *buySellOrders) (price float64, vol int64, err error) {
	price = math.NaN()
	vol = 0
	// buy orders start high and we just want to look at the highest (top) price
	buyorder, err := bs.buyOrders.Top()
	if err != nil {
		return
	}
	// sell orders start low
	for selliter := bs.sellOrders.iterator(); selliter.Next() == true; {
		sellorder := selliter.Value().(OrderState)
		//println("maxPriceOnSellSide buy sell price vol", buyorder.Price(), sellorder.Price(), sellorder.OrderQty())
		if sellorder.Price() <= buyorder.Price() {
			price = sellorder.Price()
			vol = vol + sellorder.OrderQty()
			//println("vol", vol)
		}
	}
	println("ret vol", vol)
	return
}

func minPriceOnBuySide(bs *buySellOrders) (price float64, vol int64, err error) {
	price = math.NaN()
	// sell orders start low and we just want to look at the lowest (top) price
	sellorder, err := bs.sellOrders.Top()
	if err != nil {
		return
	}
	// buy orders start high
	for buyiter := bs.buyOrders.iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(OrderState)
		//println("minPriceOnBuySide buy sell price", buyorder.Price(), sellorder.Price())
		if buyorder.Price() >= sellorder.Price() {
			price = buyorder.Price()
			vol += buyorder.LeavesQty()
		}
	}
	return
}
