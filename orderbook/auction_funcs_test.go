package orderbook

import (
	"github.com/andres-erbsen/clock"
	"github.com/shopspring/decimal"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_Auction_MinOrder(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	priceb, volb, err := minPriceOnBuySide(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.02, priceb, "min price")
	assert.AssertEqualTint64(t, 640, volb, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	err = state.minPriceBuy(bk.auctionBookOrders())
	assert.AssertTrueT(t, decimal.NewFromFloat(1.02).Equal(state.state().buy.pricelimit), "min price")
	assert.AssertEqualTint64(t, 640, state.state().buy.vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_MaxOrder(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	prices, vols, err := maxPriceOnSellSide(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.05, prices, "max sell price")
	assert.AssertEqualTint64(t, 780, vols, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	err = state.maxPriceSell(bk.auctionBookOrders())
	assert.AssertTrueT(t, decimal.NewFromFloat(1.05).Equal(state.state().sell.pricelimit), "min price")
	assert.AssertEqualTint64(t, 780, state.state().sell.vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_MaxrderVolume(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	vol, pricemin, pricemax, err := volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.02, pricemin, "min price")
	assert.AssertEqualT(t, 1.05, pricemax, "max price")
	assert.AssertEqualTint64(t, 640, vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	err = state.volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	assert.AssertTrueT(t, decimal.NewFromFloat(1.02).Equal(state.state().buy.pricelimit), "min price")
	assert.AssertTrueT(t, decimal.NewFromFloat(1.05).Equal(state.state().sell.pricelimit), "max price")
	assert.AssertEqualTint64(t, 640, state.state().clearingvol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_BuyVWAP(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	vol, _, _, err := volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	price := buyVWAP(bk.auctionBookOrders().buyOrders, vol)
	assert.AssertEqualTfloat64(t, 1.03469, price, 0.0001, "buy vwap price")
	assert.AssertEqualTint64(t, 640, vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	state.buyVWAP(bk.auctionBookOrders().buyOrders, vol)
	println("buy vwap ", state.state().buy.vwap.String())
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.03469), state.state().buy.vwap, 0.0001, "buy vwap price")
	assert.AssertEqualTint64(t, 640, state.state().clearingvol, "clearing vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_SellVWAP(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	vol, _, _, err := volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	price := sellVWAP(bk.auctionBookOrders().sellOrders, vol)
	assert.AssertEqualTfloat64(t, 1.01984, price, 0.0001, "sell vwap price")
	assert.AssertEqualTint64(t, 640, vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	state.sellVWAP(bk.auctionBookOrders().sellOrders, vol)
	println("sell vwap ", state.state().sell.vwap.String())
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.01984), state.state().sell.vwap, 0.0001, "sell vwap price")
	assert.AssertEqualTint64(t, 640, state.state().clearingvol, "clearing vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_ClearingPrice(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	price, _, err := calcClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualTfloat64(t, 1.02726, price, 0.0001, "min price")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	state.calcClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.02726), state.state().midclearingprice, 0.0001, "mid clearing price")
}

func Test_OrderBook_Auction_ClearingPricePercentages(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	price, _, _ := calcClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualTfloat64(t, 1.02726, price, 0.0001, "min price")
	buyperc, sellperc, lower, upper := calcClearingPricePercentages(price)
	println("lower upper", lower.String(), upper.String())
	toFloat := func(d decimal.Decimal) float64 {
		ret, _ := d.Float64()
		return ret
	}
	assert.AssertEqualTfloat64(t, 1.02, toFloat(lower), 0.001, "lower price")
	assert.AssertEqualTfloat64(t, 1.03, toFloat(upper), 0.001, "lower price")
	assert.AssertEqualTfloat64(t, 0.2734375, toFloat(buyperc), 0.0000001, "but perc price")
	assert.AssertEqualTfloat64(t, 0.7265625, toFloat(sellperc), 0.0000001, "sell perc price")
	state := newAuctionCloseCalculator()
	state.calcClearingPrice(bk.auctionBookOrders())
	state.calcClearingPricePercentages()
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.02726), state.state().midclearingprice, 0.0001, "mid clearing price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.02), state.state().sell.clearingprice, 0.001, "lower sell price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.03), state.state().buy.clearingprice, 0.001, "lower buy price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(0.2734375), state.state().buy.percent, 0.0000001, "buy perc price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(0.7265625), state.state().sell.percent, 0.0000001, "sell perc price")

}

func Test_OrderBook_Auction_FillOrdersStateless(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	execs, _, _, err := fillAuctionAtClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualT(t, 16, len(execs), "16")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_FillOrdersWithState(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	printOrders(bk.auctionBookOrders())
	state := newAuctionCloseCalculator()
	execs, err := state.fillAuctionAtClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualT(t, 44, len(execs), "44")
	assert.AssertEqualT(t, nil, err, "min err")
	printExecsAndOrders(execs, bk.auctionBookOrders())

	assert.Fail(t, "add more tests")
}

/*
Order id	price	vol		price	vol	order id
b1	1.05	100				1.01	300	s1
b2	1.05	50				1.01	50	s2
b3	1.05	30				1.02	100	s3
b4	1.04	100				1.03	40	s4
b5	1.03	200				1.04	250	s5
b6	1.02	60				1.05	40	s6
b7	1.02	100				1.06	200	s7
b8	0.99	200
*/
func makeOrderBook_for_OrderBook_Auction(t *testing.T) OrderBook {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())

	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli1", "b1", SideBuy, 100, 1.05))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli2", "b2", SideBuy, 50, 1.05))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli3", "b3", SideBuy, 30, 1.05))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli4", "b4", SideBuy, 100, 1.04))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli5", "b5", SideBuy, 200, 1.03))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli6", "b6", SideBuy, 60, 1.02))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli7", "b7", SideBuy, 100, 1.02))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli8", "b8", SideBuy, 200, 0.99))

	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli11", "s1", SideSell, 300, 1.01))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli21", "s2", SideSell, 50, 1.01))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli31", "s3", SideSell, 100, 1.02))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli41", "s4", SideSell, 40, 1.03))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli51", "s5", SideSell, 250, 1.04))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli61", "s6", SideSell, 40, 1.05))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli71", "s7", SideSell, 200, 1.06))

	assert.AssertEqualT(t, 8, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 7, len(bk.SellAuctionOrders()), "sell orders")

	return bk
}
