package orderbook

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func Test_OrderBook_State_TradingOpen_Closed_Market(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideSell, 101, aclock))
	e21, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id22", fixmodel.SideSell, 100, 1.05, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "fill", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "partfill", 100, 1.01, 1)

	execs, err := bk.CloseTrading()
	assert.AssertNilT(t, err, "closed ok")
	assert.AssertEqualT(t, 0, len(execs), "execs 0")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
}

func Test_OrderBook_State_TradingOpen_Closed_GoodForDay(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideSell, 101, aclock))
	e21, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id22", fixmodel.SideSell, 100, 1.05, fixmodel.TimeInForceDay, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e21), "e21 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "fill", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "partfill", 100, 1.01, 1)

	execs, err := bk.CloseTrading()
	assert.AssertNilT(t, err, "closed ok")
	assert.AssertEqualT(t, 1, len(execs), "execs 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")

	//printExecs(execs)
	fixmodel.ContainsExec(t, execs, "cli2", "id22", fixmodel.OrdStatusCanceled, "cancel", 0, 0, 1)
}
