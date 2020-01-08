package orderbookex

import (
	"fmt"
	"github.com/shopspring/decimal"
	"orderbook/obmath"
)

type auctionCloseCalculator interface {
	maxPriceSell(bs *buySellOrders) (err error)
	minPriceBuy(bs *buySellOrders) (err error)
	volBetweenPriceMinAndPriceMax(bs *buySellOrders) (err error)
	buyVWAP(orders OrderList, maxvol int64)
	sellVWAP(orders OrderList, maxvol int64)
	calcClearingPrice(bk *buySellOrders)
	calcClearingPricePercentages()
	fillAuctionAtClearingPrice(bk *buySellOrders) (execs []ExecutionReport, err error)
	state() *auctionclose
}

func newAuctionCloseCalculator() auctionCloseCalculator {
	return auctionCloseCalculator(&auctionclose{})
}

type auctionstateside struct {
	pricelimit    decimal.Decimal // min for buy max for sell
	vol           int64
	vwap          decimal.Decimal
	percent       decimal.Decimal
	clearingprice decimal.Decimal
}

type auctionclose struct {
	buy              auctionstateside
	sell             auctionstateside
	clearingvol      int64
	midclearingprice decimal.Decimal
}

func (s *auctionclose) state() *auctionclose {
	return s
}

func (s *auctionclose) maxPriceSell(bs *buySellOrders) (err error) {
	price, vol, err := maxPriceOnSellSide(bs)
	s.sell.pricelimit = decimal.NewFromFloat(price)
	s.sell.vol = vol
	return
}

func (s *auctionclose) minPriceBuy(bs *buySellOrders) (err error) {
	price, vol, err := minPriceOnBuySide(bs)
	s.buy.pricelimit = decimal.NewFromFloat(price)
	s.buy.vol = vol
	return
}

func (s *auctionclose) volBetweenPriceMinAndPriceMax(bs *buySellOrders) (err error) {
	maxvol, pricemin, pricemax, err := volBetweenPriceMinAndPriceMax(bs)
	s.buy.pricelimit = decimal.NewFromFloat(pricemin)
	s.sell.pricelimit = decimal.NewFromFloat(pricemax)
	s.clearingvol = maxvol
	return err
}

func (s *auctionclose) buyVWAP(orders OrderList, maxvol int64) {
	matchBuy := func(neworder OrderState, bookorder OrderState) bool {
		return (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() >= bookorder.Price()
	}
	s.clearingvol = maxvol
	s.buy.vwap = decimal.NewFromFloat(calcVWAP(orders, s.clearingvol, matchBuy))
}

func (s *auctionclose) sellVWAP(orders OrderList, maxvol int64) {
	matchSell := func(neworder OrderState, bookorder OrderState) bool {
		return (neworder.OrderType() == OrderTypeMarket || bookorder.OrderType() == OrderTypeMarket) || neworder.Price() <= bookorder.Price()
	}
	s.clearingvol = maxvol
	s.sell.vwap = decimal.NewFromFloat(calcVWAP(orders, s.clearingvol, matchSell))
}

func (s *auctionclose) calcClearingPrice(bk *buySellOrders) {
	s.volBetweenPriceMinAndPriceMax(bk)
	s.sellVWAP(bk.sellOrders, s.clearingvol)
	s.buyVWAP(bk.buyOrders, s.clearingvol)
	s.midclearingprice = s.sell.vwap.Add(s.buy.vwap).Div(decimal.NewFromFloat(2.0))
}

func (s *auctionclose) calcClearingPricePercentages() {
	s.buy.clearingprice = obmath.Roundup(s.midclearingprice, 2)    // upper buy price
	s.sell.clearingprice = obmath.Rounddown(s.midclearingprice, 2) // lower sell price
	percentd := s.midclearingprice.Sub(s.sell.clearingprice).Mul(decimal.NewFromFloat32(10).Pow(decimal.New(int64(2), 0)))
	s.sell.percent = percentd
	s.buy.percent = decimal.New(1, 0).Sub(percentd)
}

func (s *auctionclose) fillAuctionAtClearingPrice(bk *buySellOrders) (execs []ExecutionReport, err error) {
	execs = []ExecutionReport{}
	s.calcClearingPrice(bk)
	s.calcClearingPricePercentages()
	lowercp, _ := s.sell.clearingprice.Float64()
	uppercp, _ := s.buy.clearingprice.Float64()
	midclearingprice, _ := s.midclearingprice.Float64()
	fmt.Printf("fillAuctionAtClearingPrice (auctionclose) \n")
	fmt.Printf("buy%% | sell%% | sell limit | buy limit | clear vol | mid clearing price | upper cp | lower cp |\n")
	fmt.Printf("%v%% | %v%% | %v | %v | %d | %v | %v | %v |\n",
		s.buy.percent, s.sell.percent, s.sell.pricelimit, s.buy.pricelimit, s.clearingvol, midclearingprice, uppercp, lowercp)

	filledSellBookOrders := []OrderState{}
	filledBuyBookOrders := []OrderState{}
	var cumqty int64

	for iter := bk.buyOrders.iterator(); iter.Next() == true; {
		buyorder := iter.Value().(OrderState)
		//fmt.Printf("fillAuctionAtClearingPrice (auctionclose): buyorder %v \n", buyorder)
		for iter := bk.sellOrders.iterator(); iter.Next() == true; {
			sellorder := iter.Value().(OrderState)
			//fmt.Printf("fillAuctionAtClearingPrice (auctionclose): sellorder %v buyorder %v \n", sellorder, buyorder)
			/*
				match buy greater or equal than lower price
				match sell less or equal than upper price
			*/
			if buyorder.LeavesQty() > 0 && (buyorder.Price() >= midclearingprice && sellorder.Price() <= midclearingprice) {
				tofill := obmath.Min(sellorder.LeavesQty(), buyorder.LeavesQty())
				//fmt.Printf("to fill %d\n", tofill)
				if tofill > 0 && cumqty <= s.clearingvol {
					cumqty += tofill
					tofilld := decimal.New(tofill, 0)
					tofillbuy := tofilld.Mul(s.buy.percent).Round(0).IntPart()
					tofillsell := tofill - tofillbuy
					//fmt.Printf("fillAuctionAtClearingPrice (auctionclose) sell orders: perc %v to fill %d buy fill %d, sell fill %d\n", s.buy.percent, tofill, tofillbuy, tofillsell)
					if tofillbuy > 0 {
						buyorder.fill(tofillbuy)
						sellorder.fill(tofillbuy)
						execs = append(execs, MakeFillExecutionReport(sellorder, lowercp, tofillbuy))
						execs = append(execs, MakeFillExecutionReport(buyorder, lowercp, tofillbuy))
					}
					if tofillsell > 0 {
						buyorder.fill(tofillsell)
						sellorder.fill(tofillsell)
						execs = append(execs, MakeFillExecutionReport(sellorder, uppercp, tofillsell))
						execs = append(execs, MakeFillExecutionReport(buyorder, uppercp, tofillsell))
					}
					if sellorder.LeavesQty() == 0 {
						filledSellBookOrders = append(filledSellBookOrders, sellorder)
					}
					if buyorder.LeavesQty() == 0 {
						filledBuyBookOrders = append(filledBuyBookOrders, buyorder)
					}
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
	return
}

func maxPriceOnSellSide(bs *buySellOrders) (price float64, vol int64, err error) {
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
		}
	}
	return
}

func minPriceOnBuySide(bs *buySellOrders) (price float64, vol int64, err error) {
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

func volBetweenPriceMinAndPriceMax(bs *buySellOrders) (maxvol int64, pricemin float64, pricemax float64, err error) {
	pricemin, volb, errmin := minPriceOnBuySide(bs)
	pricemax, vols, errmax := maxPriceOnSellSide(bs)
	maxvol = obmath.Min(volb, vols)
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
			tofill := obmath.Min(bookorder.LeavesQty(), order.LeavesQty())
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
