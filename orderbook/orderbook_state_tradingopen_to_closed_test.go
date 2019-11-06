package orderbook

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_State_TradingOpen_Closed_Market(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 101, aclock))
	e21, _ := bk.NewOrder(makeLimitOrder("cli2", "id22", SideSell, 100, 1.05, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "fill", 100, 1.01, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "partfill", 100, 1.01, 1)

	execs, err := bk.CloseTrading()
	assert.AssertNilT(t, err, "closed ok")
	assert.AssertEqualT(t, 0, len(execs), "execs 0")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
}

func Test_OrderBook_State_TradingOpen_Closed_GoodForDay(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 101, aclock))
	e21, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id22", SideSell, 100, 1.05, TimeInForceDay, makeTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "fill", 100, 1.01, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "partfill", 100, 1.01, 1)

	execs, err := bk.CloseTrading()
	assert.AssertNilT(t, err, "closed ok")
	assert.AssertEqualT(t, 1, len(execs), "execs 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")

	//printExecs(execs)
	containsExec(t, execs, "cli2", "id22", OrdStatusCanceled, "cancel", 0, 0, 1)
}
