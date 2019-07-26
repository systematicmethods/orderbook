package orderbook

import (
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_State_TradingOpen_Closed_Market(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 101))
	e21, _ := bk.NewOrder(makeLimitOrder("cli2", "id22", SideSell, 100, 1.05))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "fill", 100, 1.01)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "partfill", 100, 1.01)

	execs, err := bk.CloseTrading()
	assert.AssertNilT(t, err, "closed ok")
	assert.AssertEqualT(t, 0, len(execs), "execs 0")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
}

func Test_OrderBook_State_TradingOpen_Closed_GoodForDay(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 101))
	e21, _ := bk.NewOrder(makeLimitOrderDay("cli2", "id22", SideSell, 100, 1.05))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "fill", 100, 1.01)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "partfill", 100, 1.01)

	execs, err := bk.CloseTrading()
	assert.AssertNilT(t, err, "closed ok")
	assert.AssertEqualT(t, 1, len(execs), "execs 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
}
