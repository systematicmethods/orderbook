package orderbook

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

type OrderBookAuction interface {
	OpenAuction() error
	CloseAuction() (execs []ExecutionReport, clearingPrice float64, clearingVol int64, err error)
	BuyAuctionSize() int
	SellAuctionSize() int
	BuyAuctionOrders() []OrderState
	SellAuctionOrders() []OrderState

	auctionBookOrders() *buySellOrders
}

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
func (b *orderbook) CloseAuction() (execs []ExecutionReport, clearingPrice float64, clearingVol int64, err error) {
	execs = []ExecutionReport{}
	// fillAuctionAtClearingPrice matchAuctionOrdersOnBook
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseAuction)
	if err == nil {
		var exs []ExecutionReport
		exs, clearingPrice, clearingVol, err = fillAuctionAtClearingPrice(&b.auctionOrders)
		execs = append(execs, exs...)
		exs = cancelOrders(&b.auctionOrders)
		execs = append(execs, exs...)
	}
	return
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
					tofill := min(sellorder.LeavesQty(), buyorder.LeavesQty())
					var price float64
					if buyorder.OrderType() == OrderTypeMarket {
						price = sellorder.Price()
					} else if sellorder.OrderType() == OrderTypeMarket {
						price = buyorder.Price()
					} else {
						price = buyorder.Price()
					}
					if tofill > 0 {
						if buyorder.fill(tofill) {
							bs.buyOrders.RemoveByID(buyorder.OrderID())
						}
						execs = append(execs, MakeFillExecutionReport(buyorder, price, tofill))
						if sellorder.fill(tofill) {
							bs.sellOrders.RemoveByID(sellorder.OrderID())
						}
						execs = append(execs, MakeFillExecutionReport(sellorder, price, tofill))
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

func calcClearingPrice(bk *buySellOrders) (clearingPrice float64, vol int64, err error) {
	vol, _, _, err = volBetweenPriceMinAndPriceMax(bk)
	sellvwwap := sellVWAP(bk.sellOrders, vol)
	buyvwap := buyVWAP(bk.buyOrders, vol)
	return (sellvwwap + buyvwap) / 2.0, vol, err
}

func roundup(num decimal.Decimal, places int32) decimal.Decimal {
	// math.Floor(x*100)/100
	factor := decimal.NewFromFloat32(10).Pow(decimal.New(int64(places), 0))
	return num.Mul(factor).Ceil().Div(factor)
}

func rounddown(num decimal.Decimal, places int32) decimal.Decimal {
	// math.Floor(x*100)/100
	factor := decimal.NewFromFloat32(10).Pow(decimal.New(int64(places), 0))
	return num.Mul(factor).Floor().Div(factor)
}

func calcClearingPricePercentages(clearingPrice float64) (buypercent float64, sellpercent float64, lower float64, upper float64) {
	cp := decimal.NewFromFloat(clearingPrice)
	upper, _ = roundup(cp, 2).Float64()
	lowerd := rounddown(cp, 2)
	lower, _ = rounddown(cp, 2).Float64()
	percentd := cp.Sub(lowerd).Mul(decimal.NewFromFloat32(10).Pow(decimal.New(int64(2), 0)))
	sellpercent, _ = percentd.Float64()
	buypercent, _ = decimal.New(1, 0).Sub(percentd).Float64()
	return
}

func fillAuctionAtClearingPrice(bk *buySellOrders) (execs []ExecutionReport, clearingPrice float64, clearingVol int64, err error) {
	match := func(neworder OrderState, bookorder OrderState) bool {
		return neworder.LeavesQty() >= 0
	}
	execs = []ExecutionReport{}
	clearingPrice, clearingVol, err = calcClearingPrice(bk)
	buyperc, sellperc, lowerprice, upperprice := calcClearingPricePercentages(clearingPrice)
	fmt.Printf("percent buy %v sell %v lower %v upper %v clear vol %d price %v \n", buyperc, sellperc, lowerprice, upperprice, clearingVol, clearingPrice)

	filledSellBookOrders := []OrderState{}

	buyorder := makeMarketOrderForAuction(clearingVol, clearingPrice, SideBuy)
	for iter := bk.sellOrders.iterator(); iter.Next() == true; {
		bookorder := iter.Value().(OrderState)
		if match(buyorder, bookorder) && bookorder.Price() <= buyorder.Price() {
			tofill := min(bookorder.LeavesQty(), buyorder.LeavesQty())
			if tofill > 0 {
				tofillbuy := int64(float64(tofill) * buyperc)
				tofillsell := tofill - tofillbuy
				fmt.Printf("sell orders: perc %v to fill %d buy fill %d, sell fill %d\n", buyperc, tofill, tofillbuy, tofillsell)
				buyorder.fill(tofillbuy)
				bookorder.fill(tofillbuy)
				execs = append(execs, MakeFillExecutionReport(bookorder, lowerprice, tofillbuy))
				buyorder.fill(tofillsell)
				bookorder.fill(tofillsell)
				execs = append(execs, MakeFillExecutionReport(bookorder, upperprice, tofillsell))
				if bookorder.LeavesQty() == 0 {
					filledSellBookOrders = append(filledSellBookOrders, bookorder)
				}
			}
		}
	}

	filledBuyBookOrders := []OrderState{}

	sellorder := makeMarketOrderForAuction(clearingVol, clearingPrice, SideSell)
	for iter := bk.buyOrders.iterator(); iter.Next() == true; {
		bookorder := iter.Value().(OrderState)
		if match(sellorder, bookorder) && bookorder.Price() >= sellorder.Price() {
			tofill := min(bookorder.LeavesQty(), sellorder.LeavesQty())
			if tofill > 0 {
				tofillbuy := int64(float64(tofill) * buyperc)
				tofillsell := tofill - tofillbuy
				fmt.Printf("buy orders: perc %v to fill %d buy fill %d, sell fill %d\n", buyperc, tofill, tofillbuy, tofillsell)
				sellorder.fill(tofillbuy)
				bookorder.fill(tofillbuy)
				execs = append(execs, MakeFillExecutionReport(bookorder, lowerprice, tofillbuy))
				sellorder.fill(tofillsell)
				bookorder.fill(tofillsell)
				execs = append(execs, MakeFillExecutionReport(bookorder, upperprice, tofillsell))
				if bookorder.LeavesQty() == 0 {
					filledBuyBookOrders = append(filledBuyBookOrders, bookorder)
				}
			}
		}
	}

	for _, v := range filledSellBookOrders {
		bk.sellOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBuyBookOrders {
		bk.buyOrders.RemoveByID(v.OrderID())
	}
	//printExecs2(execs)
	return
}

func volBetweenPriceMinAndPriceMax(bs *buySellOrders) (maxvol int64, pricemin float64, pricemax float64, err error) {
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

func calcVWAP(orders OrderList, maxvol int64, match func(OrderState, OrderState) bool) (vwap float64) {
	copyOfOrders := orders.orderList()
	execs := []ExecutionReport{}
	order := makeMarketOrderForAuction(maxvol, 0, SideBuy)
	for iter := copyOfOrders.iterator(); iter.Next() == true; {
		bookorder := iter.Value().(OrderState)
		if match(order, bookorder) {
			tofill := min(bookorder.LeavesQty(), order.LeavesQty())
			price := bookorder.Price()
			if tofill > 0 {
				order.fill(tofill)
				execs = append(execs, MakeFillExecutionReport(order, price, tofill))
			}
		}
	}
	return cummulativeVwapCalc(execs, maxvol)
}

func cummulativeVwapCalc(execs []ExecutionReport, maxvol int64) float64 {
	var priceXvol float64
	for _, s := range execs {
		priceXvol += s.LastPrice() * float64(s.LastQty()) / float64(maxvol)
	}
	return priceXvol
}

func buyVWAP(orders OrderList, maxvol int64) (vwap float64) {
	matchBuy := func(neworder OrderState, bookorder OrderState) bool {
		return (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() >= bookorder.Price()
	}
	return calcVWAP(orders, maxvol, matchBuy)
}

func sellVWAP(orders OrderList, maxvol int64) (vwap float64) {
	matchSell := func(neworder OrderState, bookorder OrderState) bool {
		return (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() <= bookorder.Price()
	}
	return calcVWAP(orders, maxvol, matchSell)
}

func printExecs2(execs []ExecutionReport) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
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
	//println("ret vol", vol)
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

func makeMarketOrderForAuction(qty int64, price float64, side Side) OrderState {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	ordertype := OrderTypeMarket
	if price != 0 {
		ordertype = OrderTypeLimit
	}
	return MakeOrderState(
		"",
		"",
		"",
		side,
		price,
		qty,
		ordertype,
		TimeInForceGoodForAuction,
		dt,
		dt,
		dt,
		dt,
		"",
		uuid.New(),
		qty,
		0,
		OrdStatusNew,
	)
}