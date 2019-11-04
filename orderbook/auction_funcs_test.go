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
	assert.AssertEqualT(t, int64(640), volb, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	err = state.minPriceBuy(bk.auctionBookOrders())
	assert.AssertTrueT(t, decimal.NewFromFloat(1.02).Equal(state.state().buy.pricelimit), "min price")
	assert.AssertEqualT(t, int64(640), state.state().buy.vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_MaxOrder(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	prices, vols, err := maxPriceOnSellSide(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.05, prices, "max sell price")
	assert.AssertEqualT(t, int64(780), vols, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	err = state.maxPriceSell(bk.auctionBookOrders())
	assert.AssertTrueT(t, decimal.NewFromFloat(1.05).Equal(state.state().sell.pricelimit), "min price")
	assert.AssertEqualT(t, int64(780), state.state().sell.vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_MaxrderVolume(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	vol, pricemin, pricemax, err := volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.02, pricemin, "min price")
	assert.AssertEqualT(t, 1.05, pricemax, "max price")
	assert.AssertEqualT(t, int64(640), vol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
	state := newAuctionCloseCalculator()
	err = state.volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	assert.AssertTrueT(t, decimal.NewFromFloat(1.02).Equal(state.state().buy.pricelimit), "min price")
	assert.AssertTrueT(t, decimal.NewFromFloat(1.05).Equal(state.state().sell.pricelimit), "max price")
	assert.AssertEqualT(t, int64(640), state.state().clearingvol, "buy vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_BuyVWAP(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	state := newAuctionCloseCalculator()
	err := state.volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	state.buyVWAP(bk.auctionBookOrders().buyOrders, state.state().clearingvol)
	println("buy vwap ", state.state().buy.vwap.String())
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.03469), state.state().buy.vwap, 0.0001, "buy vwap price")
	assert.AssertEqualT(t, int64(640), state.state().clearingvol, "clearing vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_SellVWAP(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	state := newAuctionCloseCalculator()
	err := state.volBetweenPriceMinAndPriceMax(bk.auctionBookOrders())
	state.sellVWAP(bk.auctionBookOrders().sellOrders, state.state().clearingvol)
	println("sell vwap ", state.state().sell.vwap.String())
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.01984), state.state().sell.vwap, 0.0001, "sell vwap price")
	assert.AssertEqualT(t, int64(640), state.state().clearingvol, "clearing vol")
	assert.AssertEqualT(t, nil, err, "min err")
}

func Test_OrderBook_Auction_ClearingPrice(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	state := newAuctionCloseCalculator()
	state.calcClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.02726), state.state().midclearingprice, 0.0001, "mid clearing price")
}

func Test_OrderBook_Auction_ClearingPricePercentages(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	state := newAuctionCloseCalculator()
	state.calcClearingPrice(bk.auctionBookOrders())
	state.calcClearingPricePercentages()
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.02726), state.state().midclearingprice, 0.0001, "mid clearing price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.02), state.state().sell.clearingprice, 0.001, "lower sell price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(1.03), state.state().buy.clearingprice, 0.001, "lower buy price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(0.2734375), state.state().buy.percent, 0.0000001, "buy perc price")
	assert.AssertEqualTdecimal(t, decimal.NewFromFloat(0.7265625), state.state().sell.percent, 0.0000001, "sell perc price")

}

func Test_OrderBook_Auction_FillOrdersWithState(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	//printOrders(bk.auctionBookOrders())
	//buyorders := bk.auctionBookOrders().buyOrders.Orders()
	//sellorders := bk.auctionBookOrders().sellOrders.Orders()
	state := newAuctionCloseCalculator()
	execs, err := state.fillAuctionAtClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualT(t, 28, len(execs), "28")
	assert.AssertEqualT(t, nil, err, "min err")

	csv := `id|clientid|clordid|side|lastprice|lastqty|status|price|qty|ordstatus
e0|cli11|s1|sell|1.02|27|PartiallyFilled|1.01|300|Filled
e1|cli1|b1|buy|1.02|27|PartiallyFilled|1.05|100|Filled
e2|cli11|s1|sell|1.03|73|PartiallyFilled|1.01|300|Filled
e3|cli1|b1|buy|1.03|73|Filled|1.05|100|Filled
e4|cli11|s1|sell|1.02|14|PartiallyFilled|1.01|300|Filled
e5|cli2|b2|buy|1.02|14|PartiallyFilled|1.05|50|Filled
e6|cli11|s1|sell|1.03|36|PartiallyFilled|1.01|300|Filled
e7|cli2|b2|buy|1.03|36|Filled|1.05|50|Filled
e8|cli11|s1|sell|1.02|8|PartiallyFilled|1.01|300|Filled
e9|cli3|b3|buy|1.02|8|PartiallyFilled|1.05|30|Filled
e10|cli11|s1|sell|1.03|22|PartiallyFilled|1.01|300|Filled
e11|cli3|b3|buy|1.03|22|Filled|1.05|30|Filled
e12|cli11|s1|sell|1.02|27|PartiallyFilled|1.01|300|Filled
e13|cli4|b4|buy|1.02|27|PartiallyFilled|1.04|100|Filled
e14|cli11|s1|sell|1.03|73|PartiallyFilled|1.01|300|Filled
e15|cli4|b4|buy|1.03|73|Filled|1.04|100|Filled
e16|cli11|s1|sell|1.02|5|PartiallyFilled|1.01|300|Filled
e17|cli5|b5|buy|1.02|5|PartiallyFilled|1.03|200|PartiallyFilled
e18|cli11|s1|sell|1.03|15|Filled|1.01|300|Filled
e19|cli5|b5|buy|1.03|15|PartiallyFilled|1.03|200|PartiallyFilled
e20|cli21|s2|sell|1.02|14|PartiallyFilled|1.01|50|Filled
e21|cli5|b5|buy|1.02|14|PartiallyFilled|1.03|200|PartiallyFilled
e22|cli21|s2|sell|1.03|36|Filled|1.01|50|Filled
e23|cli5|b5|buy|1.03|36|PartiallyFilled|1.03|200|PartiallyFilled
e24|cli31|s3|sell|1.02|27|PartiallyFilled|1.02|100|Filled
e25|cli5|b5|buy|1.02|27|PartiallyFilled|1.03|200|PartiallyFilled
e26|cli31|s3|sell|1.03|73|Filled|1.02|100|Filled
e27|cli5|b5|buy|1.03|73|PartiallyFilled|1.03|200|PartiallyFilled
`
	expected := loadExecCSV(csv)

	for _, ex := range expected {
		containsExecCSV(t, execs, ex, "execs")
	}
	//printExecsAndOrders(execs, bk.auctionBookOrders(), buyorders, sellorders)
}

func Test_OrderBook_Auction_FillOrdersWithStateExterme(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction_Extreme(t)
	//printOrders(bk.auctionBookOrders())
	//buyorders := bk.auctionBookOrders().buyOrders.Orders()
	//sellorders := bk.auctionBookOrders().sellOrders.Orders()
	state := newAuctionCloseCalculator()
	execs, err := state.fillAuctionAtClearingPrice(bk.auctionBookOrders())
	assert.AssertEqualT(t, 16, len(execs), "16")
	assert.AssertEqualT(t, nil, err, "min err")
	csv := `id|clientid|clordid|side|lastprice|lastqty|status|price|qty|ordstatus
e0|cli21|s2|sell|1.76|19|PartiallyFilled|1.01|50|Filled
e1|cli3|b3|buy|1.76|19|PartiallyFilled|2.1|30|Filled
e2|cli21|s2|sell|1.77|11|PartiallyFilled|1.01|50|Filled
e3|cli3|b3|buy|1.77|11|Filled|2.1|30|Filled
e4|cli21|s2|sell|1.76|13|PartiallyFilled|1.01|50|Filled
e5|cli5|b5|buy|1.76|13|PartiallyFilled|2.06|200|PartiallyFilled
e6|cli21|s2|sell|1.77|7|Filled|1.01|50|Filled
e7|cli5|b5|buy|1.77|7|PartiallyFilled|2.06|200|PartiallyFilled
e8|cli41|s4|sell|1.76|25|PartiallyFilled|1.03|40|Filled
e9|cli5|b5|buy|1.76|25|PartiallyFilled|2.06|200|PartiallyFilled
e10|cli41|s4|sell|1.77|15|Filled|1.03|40|Filled
e11|cli5|b5|buy|1.77|15|PartiallyFilled|2.06|200|PartiallyFilled
e12|cli61|s6|sell|1.76|25|PartiallyFilled|1.05|40|Filled
e13|cli5|b5|buy|1.76|25|PartiallyFilled|2.06|200|PartiallyFilled
e14|cli61|s6|sell|1.77|15|Filled|1.05|40|Filled
e15|cli5|b5|buy|1.77|15|PartiallyFilled|2.06|200|PartiallyFilled
`
	expected := loadExecCSV(csv)

	for _, ex := range expected {
		containsExecCSV(t, execs, ex, "execs")
	}
	//printExecsAndOrders(execs, bk.auctionBookOrders(), buyorders, sellorders)
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

func makeOrderBook_for_OrderBook_Auction_Extreme(t *testing.T) OrderBook {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())

	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli1", "b1", SideBuy, 100, 1.05/4))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli2", "b2", SideBuy, 50, 1.05))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli3", "b3", SideBuy, 30, 1.05*2))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli4", "b4", SideBuy, 100, 1.04))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli5", "b5", SideBuy, 200, 1.03*2))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli6", "b6", SideBuy, 60, 1.02))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli7", "b7", SideBuy, 100, 1.02*2))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli8", "b8", SideBuy, 200, 0.99))

	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli11", "s1", SideSell, 300, 1.01*4))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli21", "s2", SideSell, 50, 1.01))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli31", "s3", SideSell, 100, 1.02*2))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli41", "s4", SideSell, 40, 1.03))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli51", "s5", SideSell, 250, 1.04*2))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli61", "s6", SideSell, 40, 1.05))
	_, _ = bk.NewOrder(makeAuctionLimitOrder("cli71", "s7", SideSell, 200, 1.06*2))

	assert.AssertEqualT(t, 8, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 7, len(bk.SellAuctionOrders()), "sell orders")

	return bk
}
