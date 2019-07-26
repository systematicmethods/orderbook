package orderbook

import (
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_Auction_AddBuySellOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenAuction)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 101, 1.03))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuyAuctionSize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellAuctionSize(), 1, "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0)
}

func Test_OrderBook_Auction_MatchBuySellOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenAuction)

	e10, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideBuy, 100, 1.01))
	e11, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id3", SideBuy, 100, 1.01))
	e12, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id4", SideBuy, 100, 1.01))
	e13, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id5", SideBuy, 100, 1.01))
	e21, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id21", SideSell, 101, 1.00))
	e22, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id22", SideSell, 101, 1.00))

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 1, len(e22), "e22 empty")
	assert.AssertEqualT(t, 4, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 2, len(bk.SellAuctionOrders()), "sell orders")

	e3, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 10, len(e3), "e3 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")

	//printExecs(e3)

	containsExec(t, e3, "cli1", "id21", OrdStatusPartiallyFilled, "part fill", 100, 1.01)
	containsExec(t, e3, "cli1", "id21", OrdStatusFilled, "fill", 1, 1.01)
	containsExec(t, e3, "cli1", "id22", OrdStatusPartiallyFilled, "part fill", 99, 1.01)
	containsExec(t, e3, "cli1", "id22", OrdStatusFilled, "fill", 2, 1.01)

	containsExec(t, e3, "cli2", "id2", OrdStatusFilled, "fill", 100, 1.01)
	containsExec(t, e3, "cli2", "id3", OrdStatusPartiallyFilled, "fill", 1, 1.01)
	containsExec(t, e3, "cli2", "id3", OrdStatusFilled, "fill", 99, 1.01)
	containsExec(t, e3, "cli2", "id4", OrdStatusPartiallyFilled, "fill", 2, 1.01)
	containsExec(t, e3, "cli2", "id4", OrdStatusCanceled, "cancel", 0, 0)
	containsExec(t, e3, "cli2", "id5", OrdStatusCanceled, "cancel", 0, 0)
}

func Test_OrderBook_Auction_MatchSellBuyOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenAuction)

	e10, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 100, 1.00))
	e11, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id3", SideSell, 100, 1.00))
	e12, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id4", SideSell, 100, 1.00))
	e13, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id5", SideSell, 100, 1.00))
	e21, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id21", SideBuy, 101, 1.01))
	e22, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id22", SideBuy, 101, 1.01))

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 1, len(e22), "e22 empty")
	assert.AssertEqualT(t, 2, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 4, len(bk.SellAuctionOrders()), "sell orders")

	e3, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 10, len(e3), "e3 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")

	//printExecs(e3)

	containsExec(t, e3, "cli1", "id21", OrdStatusPartiallyFilled, "part fill", 100, 1.01)
	containsExec(t, e3, "cli1", "id21", OrdStatusFilled, "fill", 1, 1.01)
	containsExec(t, e3, "cli1", "id22", OrdStatusPartiallyFilled, "part fill", 99, 1.01)
	containsExec(t, e3, "cli1", "id22", OrdStatusFilled, "fill", 2, 1.01)

	containsExec(t, e3, "cli2", "id2", OrdStatusFilled, "fill", 100, 1.01)
	containsExec(t, e3, "cli2", "id3", OrdStatusPartiallyFilled, "fill", 1, 1.01)
	containsExec(t, e3, "cli2", "id3", OrdStatusFilled, "fill", 99, 1.01)
	containsExec(t, e3, "cli2", "id4", OrdStatusPartiallyFilled, "fill", 2, 1.01)
	containsExec(t, e3, "cli2", "id4", OrdStatusCanceled, "cancel", 0, 0)
	containsExec(t, e3, "cli2", "id5", OrdStatusCanceled, "cancel", 0, 0)
}

func Test_OrderBook_Auction_MatchSellBuyOrderPlaceDuringTrading(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeNoTrading)

	e10, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 100, 1.00))
	e11, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id3", SideSell, 100, 1.00))
	e12, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id4", SideSell, 100, 1.00))
	e13, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id5", SideSell, 100, 1.00))
	e21, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id21", SideBuy, 101, 1.01))
	e22, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id22", SideBuy, 101, 1.01))

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 1, len(e22), "e22 empty")
	assert.AssertEqualT(t, 2, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 4, len(bk.SellAuctionOrders()), "sell orders")

	_, err := bk.OpenTrading()
	assert.AssertNilT(t, err, "should close trading")
	_, err = bk.CloseTrading()
	assert.AssertNilT(t, err, "should close trading")
	err = bk.OpenAuction()
	assert.AssertNilT(t, err, "should open an auction")

	e3, err := bk.CloseAuction()
	assert.AssertEqualT(t, 10, len(e3), "e3 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")

	//printExecs(e3)

	containsExec(t, e3, "cli1", "id21", OrdStatusPartiallyFilled, "part fill", 100, 1.01)
	containsExec(t, e3, "cli1", "id21", OrdStatusFilled, "fill", 1, 1.01)
	containsExec(t, e3, "cli1", "id22", OrdStatusPartiallyFilled, "part fill", 99, 1.01)
	containsExec(t, e3, "cli1", "id22", OrdStatusFilled, "fill", 2, 1.01)

	containsExec(t, e3, "cli2", "id2", OrdStatusFilled, "fill", 100, 1.01)
	containsExec(t, e3, "cli2", "id3", OrdStatusPartiallyFilled, "fill", 1, 1.01)
	containsExec(t, e3, "cli2", "id3", OrdStatusFilled, "fill", 99, 1.01)
	containsExec(t, e3, "cli2", "id4", OrdStatusPartiallyFilled, "fill", 2, 1.01)
	containsExec(t, e3, "cli2", "id4", OrdStatusCanceled, "cancel", 0, 0)
	containsExec(t, e3, "cli2", "id5", OrdStatusCanceled, "cancel", 0, 0)

}
