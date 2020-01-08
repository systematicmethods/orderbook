package auction

import (
	"fmt"
	"github.com/shopspring/decimal"
	"orderbook/fixmodel"
	"orderbook/obmath"
	"orderbook/orderstate"
)

type auctionCloseCalculator interface {
	maxPriceSell(bs *BuySellOrders) (err error)
	minPriceBuy(bs *BuySellOrders) (err error)
	volBetweenPriceMinAndPriceMax(bs *BuySellOrders) (err error)
	buyVWAP(orders *orderstate.Orderlist, maxvol int64)
	sellVWAP(orders *orderstate.Orderlist, maxvol int64)
	calcClearingPrice(bk *BuySellOrders)
	calcClearingPricePercentages()
	fillAuctionAtClearingPrice(bk *BuySellOrders) (execs []*fixmodel.ExecutionReport, err error)
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

func (s *auctionclose) maxPriceSell(bs *BuySellOrders) (err error) {
	price, vol, err := maxPriceOnSellSide(bs)
	s.sell.pricelimit = decimal.NewFromFloat(price)
	s.sell.vol = vol
	return
}

func (s *auctionclose) minPriceBuy(bs *BuySellOrders) (err error) {
	price, vol, err := minPriceOnBuySide(bs)
	s.buy.pricelimit = decimal.NewFromFloat(price)
	s.buy.vol = vol
	return
}

func (s *auctionclose) volBetweenPriceMinAndPriceMax(bs *BuySellOrders) (err error) {
	maxvol, pricemin, pricemax, err := volBetweenPriceMinAndPriceMax(bs)
	s.buy.pricelimit = decimal.NewFromFloat(pricemin)
	s.sell.pricelimit = decimal.NewFromFloat(pricemax)
	s.clearingvol = maxvol
	return err
}

func (s *auctionclose) buyVWAP(orders *orderstate.Orderlist, maxvol int64) {
	matchBuy := func(neworder *orderstate.OrderState, bookorder *orderstate.OrderState) bool {
		return (neworder.OrderType() == fixmodel.OrderTypeMarket || bookorder.OrderType() == fixmodel.OrderTypeMarket) || neworder.Price() >= bookorder.Price()
	}
	s.clearingvol = maxvol
	s.buy.vwap = decimal.NewFromFloat(calcVWAP(orders, s.clearingvol, matchBuy))
}

func (s *auctionclose) sellVWAP(orders *orderstate.Orderlist, maxvol int64) {
	matchSell := func(neworder *orderstate.OrderState, bookorder *orderstate.OrderState) bool {
		return (neworder.OrderType() == fixmodel.OrderTypeMarket || bookorder.OrderType() == fixmodel.OrderTypeMarket) || neworder.Price() <= bookorder.Price()
	}
	s.clearingvol = maxvol
	s.sell.vwap = decimal.NewFromFloat(calcVWAP(orders, s.clearingvol, matchSell))
}

func (s *auctionclose) calcClearingPrice(bk *BuySellOrders) {
	s.volBetweenPriceMinAndPriceMax(bk)
	s.sellVWAP(bk.SellOrders, s.clearingvol)
	s.buyVWAP(bk.BuyOrders, s.clearingvol)
	s.midclearingprice = s.sell.vwap.Add(s.buy.vwap).Div(decimal.NewFromFloat(2.0))
}

func (s *auctionclose) calcClearingPricePercentages() {
	s.buy.clearingprice = obmath.Roundup(s.midclearingprice, 2)    // upper buy price
	s.sell.clearingprice = obmath.Rounddown(s.midclearingprice, 2) // lower sell price
	percentd := s.midclearingprice.Sub(s.sell.clearingprice).Mul(decimal.NewFromFloat32(10).Pow(decimal.New(int64(2), 0)))
	s.sell.percent = percentd
	s.buy.percent = decimal.New(1, 0).Sub(percentd)
}

func (s *auctionclose) fillAuctionAtClearingPrice(bk *BuySellOrders) (execs []*fixmodel.ExecutionReport, err error) {
	execs = []*fixmodel.ExecutionReport{}
	s.calcClearingPrice(bk)
	s.calcClearingPricePercentages()
	lowercp, _ := s.sell.clearingprice.Float64()
	uppercp, _ := s.buy.clearingprice.Float64()
	midclearingprice, _ := s.midclearingprice.Float64()
	fmt.Printf("fillAuctionAtClearingPrice (auctionclose) \n")
	fmt.Printf("buy%% | sell%% | sell limit | buy limit | clear vol | mid clearing price | upper cp | lower cp |\n")
	fmt.Printf("%v%% | %v%% | %v | %v | %d | %v | %v | %v |\n",
		s.buy.percent, s.sell.percent, s.sell.pricelimit, s.buy.pricelimit, s.clearingvol, midclearingprice, uppercp, lowercp)

	filledSellBookOrders := []*orderstate.OrderState{}
	filledBuyBookOrders := []*orderstate.OrderState{}
	var cumqty int64

	for iter := bk.BuyOrders.Iterator(); iter.Next() == true; {
		buyorder := iter.Value().(*orderstate.OrderState)
		//fmt.Printf("fillAuctionAtClearingPrice (auctionclose): buyorder %v \n", buyorder)
		for iter := bk.SellOrders.Iterator(); iter.Next() == true; {
			sellorder := iter.Value().(*orderstate.OrderState)
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
						buyorder.Fill(tofillbuy)
						sellorder.Fill(tofillbuy)
						execs = append(execs, orderstate.NewFillExecutionReport(sellorder, lowercp, tofillbuy))
						execs = append(execs, orderstate.NewFillExecutionReport(buyorder, lowercp, tofillbuy))
					}
					if tofillsell > 0 {
						buyorder.Fill(tofillsell)
						sellorder.Fill(tofillsell)
						execs = append(execs, orderstate.NewFillExecutionReport(sellorder, uppercp, tofillsell))
						execs = append(execs, orderstate.NewFillExecutionReport(buyorder, uppercp, tofillsell))
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
		bk.SellOrders.RemoveByID(v.OrderID())
	}
	for _, v := range filledBuyBookOrders {
		bk.BuyOrders.RemoveByID(v.OrderID())
	}
	return
}

func maxPriceOnSellSide(bs *BuySellOrders) (price float64, vol int64, err error) {
	vol = 0
	// buy orders start high and we just want to look at the highest (top) price
	buyorder := bs.BuyOrders.Top()
	// sell orders start low
	for selliter := bs.SellOrders.Iterator(); selliter.Next() == true; {
		sellorder := selliter.Value().(*orderstate.OrderState)
		//println("maxPriceOnSellSide buy sell price vol", buyorder.Price(), sellorder.Price(), sellorder.OrderQty())
		if sellorder.Price() <= buyorder.Price() {
			price = sellorder.Price()
			vol = vol + sellorder.OrderQty()
		}
	}
	return
}

func minPriceOnBuySide(bs *BuySellOrders) (price float64, vol int64, err error) {
	// sell orders start low and we just want to look at the lowest (top) price
	sellorder := bs.SellOrders.Top()
	// buy orders start high
	for buyiter := bs.BuyOrders.Iterator(); buyiter.Next() == true; {
		buyorder := buyiter.Value().(*orderstate.OrderState)
		//println("minPriceOnBuySide buy sell price", buyorder.Price(), sellorder.Price())
		if buyorder.Price() >= sellorder.Price() {
			price = buyorder.Price()
			vol += buyorder.LeavesQty()
		}
	}
	return
}

func volBetweenPriceMinAndPriceMax(bs *BuySellOrders) (maxvol int64, pricemin float64, pricemax float64, err error) {
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

func calcVWAP(orders *orderstate.Orderlist, maxvol int64, match func(*orderstate.OrderState, *orderstate.OrderState) bool) (vwap float64) {
	copyOfOrders := orders.CopyList()
	execs := []*fixmodel.ExecutionReport{}
	order := NewMarketOrderForAuction(maxvol, 0, fixmodel.SideBuy)
	for iter := copyOfOrders.Iterator(); iter.Next() == true; {
		bookorder := iter.Value().(*orderstate.OrderState)
		if match(order, bookorder) {
			tofill := obmath.Min(bookorder.LeavesQty(), order.LeavesQty())
			price := bookorder.Price()
			if tofill > 0 {
				order.Fill(tofill)
				execs = append(execs, orderstate.NewFillExecutionReport(order, price, tofill))
			}
		}
	}
	return cummulativeVwapCalc(execs, maxvol)
}

func cummulativeVwapCalc(execs []*fixmodel.ExecutionReport, maxvol int64) float64 {
	var priceXvol float64
	for _, s := range execs {
		priceXvol += s.LastPrice() * float64(s.LastQty()) / float64(maxvol)
	}
	return priceXvol
}
