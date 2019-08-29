package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_Auction_MinOrder(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	priceb, volb, errb := minPriceOnBuySide(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.02, priceb, "min price")
	assert.AssertEqualTint64(t, 640, volb, "buy vol")
	assert.AssertEqualT(t, nil, errb, "min err")
}

func Test_OrderBook_Auction_MaxOrder(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	prices, vols, errs := maxPriceOnSellSide(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.05, prices, "max sell price")
	assert.AssertEqualTint64(t, 780, vols, "buy vol")
	assert.AssertEqualT(t, nil, errs, "min err")
}

func Test_OrderBook_Auction_MaxrderVolume(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	vol, pricemin, pricemax, errs := volMaxPriceMinPriceMax(bk.auctionBookOrders())
	assert.AssertEqualT(t, 1.02, pricemin, "min price")
	assert.AssertEqualT(t, 1.05, pricemax, "max price")
	assert.AssertEqualTint64(t, 640, vol, "buy vol")
	assert.AssertEqualT(t, nil, errs, "min err")
}

func Test_OrderBook_Auction_BuyVWAP(t *testing.T) {
	bk := makeOrderBook_for_OrderBook_Auction(t)
	vol, pricemin, pricemax, errs := volMaxPriceMinPriceMax(bk.auctionBookOrders())
	vwapp := buyVWAP(bk.auctionBookOrders().buyOrders, vol)
	assert.AssertEqualT(t, 1.0347, vwapp, "min price")
	assert.AssertEqualT(t, 1.02, pricemin, "min price")
	assert.AssertEqualT(t, 1.05, pricemax, "max price")
	assert.AssertEqualTint64(t, 640, vol, "buy vol")
	assert.AssertEqualT(t, nil, errs, "min err")
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
