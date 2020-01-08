package orderbookex

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_Auction_AddBuySellOrder(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 101, 1.03, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuyAuctionSize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellAuctionSize(), 1, "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new order", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0, 1)
}

func Test_OrderBook_Auction_AddBuySellOrderMid(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id1", SideBuy, 100, 1.02, aclock))
	e2, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuyAuctionSize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellAuctionSize(), 1, "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new order", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0, 1)

	e3, clrPrice, clrVol, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 5, len(e3), "e3 5")
	assert.AssertEqualTfloat64(t, 1.015, clrPrice, 0.0001, "clearing price")
	assert.AssertEqualT(t, int64(100), clrVol, "clearing vol")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")

	containsExec(t, e3, "cli1", "id1", OrdStatusPartiallyFilled, "part fill", 50, 1.01, 1)
}

func Test_OrderBook_Auction_MatchBuySellOrder(t *testing.T) {

	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e10, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideBuy, 100, 1.01, aclock))   // 2
	e11, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id3", SideBuy, 100, 1.01, aclock))   // 2
	e12, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id4", SideBuy, 100, 1.01, aclock))   // 2 cancel
	e13, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id5", SideBuy, 100, 1.01, aclock))   // cancel
	e21, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id21", SideSell, 101, 1.00, aclock)) // 2
	e22, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id22", SideSell, 101, 1.00, aclock)) // 2

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 1, len(e22), "e22 empty")
	assert.AssertEqualT(t, 4, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 2, len(bk.SellAuctionOrders()), "sell orders")

	e3, clrPrice, clrVol, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 16, len(e3), "e3 ")
	assert.AssertEqualTfloat64(t, 1.005, clrPrice, 0.0001, "clearing price")
	assert.AssertEqualT(t, int64(202), clrVol, "clearing vol")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")

	csv := `id|clientid|clordid|side|lastprice|lastqty|status
e0|cli1|id21|sell|1|50|PartiallyFilled
e1|cli2|id2|buy|1|50|PartiallyFilled
e2|cli1|id21|sell|1.01|50|PartiallyFilled
e3|cli2|id2|buy|1.01|50|Filled
e4|cli1|id21|sell|1|1|Filled
e5|cli2|id3|buy|1|1|PartiallyFilled
e6|cli1|id22|sell|1|50|PartiallyFilled
e7|cli2|id3|buy|1|50|PartiallyFilled
e8|cli1|id22|sell|1.01|49|PartiallyFilled
e9|cli2|id3|buy|1.01|49|Filled
e10|cli1|id22|sell|1|1|PartiallyFilled
e11|cli2|id4|buy|1|1|PartiallyFilled
e12|cli1|id22|sell|1.01|1|Filled
e13|cli2|id4|buy|1.01|1|PartiallyFilled
e14|cli2|id4|buy|0|0|Cancelled
e15|cli2|id5|buy|0|0|Cancelled`
	expected := loadExecCSV(csv)

	for _, ex := range expected {
		containsExecCSV(t, e3, ex, "e3")
	}

	//printExecs(e3)
}

func Test_OrderBook_Auction_MatchSellBuyOrder(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e10, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 100, 1.00, aclock)) // 2
	e11, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id3", SideSell, 100, 1.00, aclock)) // 2
	e12, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id4", SideSell, 100, 1.00, aclock)) // 2 + 1
	e13, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id5", SideSell, 100, 1.00, aclock)) // 1
	e21, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id21", SideBuy, 101, 1.01, aclock)) // 2 + 1
	e22, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id22", SideBuy, 101, 1.01, aclock)) // 2 + 1

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 1, len(e22), "e22 empty")
	assert.AssertEqualT(t, 2, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 4, len(bk.SellAuctionOrders()), "sell orders")

	e3, clrPrice, clrVol, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 16, len(e3), "e3 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")
	assert.AssertEqualT(t, 1.005, clrPrice, "clearing price")
	assert.AssertEqualT(t, int64(202), clrVol, "clearing vol")

	csv := `id|clientid|clordid|side|lastprice|lastqty|status
e0|cli2|id2|sell|1|50|PartiallyFilled
e1|cli1|id21|buy|1|50|PartiallyFilled
e2|cli2|id2|sell|1.01|50|Filled
e3|cli1|id21|buy|1.01|50|PartiallyFilled
e4|cli2|id3|sell|1|1|PartiallyFilled
e5|cli1|id21|buy|1|1|Filled
e6|cli2|id3|sell|1|50|PartiallyFilled
e7|cli1|id22|buy|1|50|PartiallyFilled
e8|cli2|id3|sell|1.01|49|Filled
e9|cli1|id22|buy|1.01|49|PartiallyFilled
e10|cli2|id4|sell|1|1|PartiallyFilled
e11|cli1|id22|buy|1|1|PartiallyFilled
e12|cli2|id4|sell|1.01|1|PartiallyFilled
e13|cli1|id22|buy|1.01|1|Filled
e14|cli2|id4|sell|0|0|Cancelled
e15|cli2|id5|sell|0|0|Cancelled`
	expected := loadExecCSV(csv)

	for _, ex := range expected {
		containsExecCSV(t, e3, ex, "e3")
	}

	//printExecs(e3)
}

func Test_OrderBook_Auction_MatchSellBuyOrderPlaceDuringTrading(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e10, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 100, 1.00, aclock)) // 2x50
	e11, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id3", SideSell, 100, 1.00, aclock)) // 1x1 1x50 1x49
	e12, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id4", SideSell, 100, 1.00, aclock)) // 1x1 1x1 1xC
	e13, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id5", SideSell, 100, 1.00, aclock)) // 1xC
	e21, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id21", SideBuy, 101, 1.01, aclock)) // 2x50 1x1
	e22, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id22", SideBuy, 101, 1.01, aclock)) // 1x50 1x49 1x1 1x1

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 1, len(e22), "e22 empty")
	assert.AssertEqualT(t, 2, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 4, len(bk.SellAuctionOrders()), "sell orders")

	_, err := bk.OpenTrading()
	assert.AssertNilT(t, err, "should open trading")
	_, err = bk.CloseTrading()
	assert.AssertNilT(t, err, "should close trading")
	err = bk.OpenAuction()
	assert.AssertNilT(t, err, "should open auction")

	e3, clrPrice, clrVol, err := bk.CloseAuction()
	assert.AssertEqualT(t, 16, len(e3), "e3 wrong")
	assert.AssertEqualT(t, 0, len(bk.BuyAuctionOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellAuctionOrders()), "sell orders")
	assert.AssertEqualT(t, 1.005, clrPrice, "clearing price")
	assert.AssertEqualT(t, int64(202), clrVol, "clearing vol")

	csv := `id|clientid|clordid|side|lastprice|lastqty|status
e0|cli2|id2|sell|1|50|PartiallyFilled
e1|cli1|id21|buy|1|50|PartiallyFilled
e2|cli2|id2|sell|1.01|50|Filled
e3|cli1|id21|buy|1.01|50|PartiallyFilled
e4|cli2|id3|sell|1|1|PartiallyFilled
e5|cli1|id21|buy|1|1|Filled
e6|cli2|id3|sell|1|50|PartiallyFilled
e7|cli1|id22|buy|1|50|PartiallyFilled
e8|cli2|id3|sell|1.01|49|Filled
e9|cli1|id22|buy|1.01|49|PartiallyFilled
e10|cli2|id4|sell|1|1|PartiallyFilled
e11|cli1|id22|buy|1|1|PartiallyFilled
e12|cli2|id4|sell|1.01|1|PartiallyFilled
e13|cli1|id22|buy|1.01|1|Filled
e14|cli2|id4|sell|0|0|Cancelled
e15|cli2|id5|sell|0|0|Cancelled`
	expected := loadExecCSV(csv)

	for _, ex := range expected {
		containsExecCSV(t, e3, ex, "e3")
	}

	//printExecs(e3)
}
