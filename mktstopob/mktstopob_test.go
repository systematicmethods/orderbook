package mktstopob

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func Test_OrderBook_Stop_Sell(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id11", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id21", fixmodel.SideSell, 100, 1.00, aclock))
	e3, _ := bk.NewOrder(makeStopOrder("cli1", "id12", fixmodel.SideSell, 100, 0.70, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1")
	assert.AssertEqualT(t, 3, len(e2), "e2")
	assert.AssertEqualT(t, 1, len(e3), "e3")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size")
	assert.AssertEqualT(t, 0, bk.BuyStopSize(), "buy stop size")
	assert.AssertEqualT(t, 1, bk.SellStopSize(), "sell stop size")

	e4, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id22", fixmodel.SideBuy, 100, 0.50, aclock))
	assert.AssertEqualT(t, 4, len(e4), "e4")

	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size")
	assert.AssertEqualT(t, 0, bk.BuyStopSize(), "buy size")
	assert.AssertEqualT(t, 0, bk.SellStopSize(), "sell size")

	fixmodel.ContainsExec(t, e4, "cli2", "id22", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e4, "cli2", "id22", fixmodel.OrdStatusFilled, "filled", 100, 0.50, 1)
	fixmodel.ContainsExec(t, e4, "cli1", "id12", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e4, "cli1", "id12", fixmodel.OrdStatusFilled, "filled", 100, 0.50, 1)
}

func Test_OrderBook_Stop_Buy(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id11", fixmodel.SideSell, 100, 1.00, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id21", fixmodel.SideBuy, 100, 1.01, aclock))
	e3, _ := bk.NewOrder(makeStopOrder("cli1", "id12", fixmodel.SideBuy, 100, 1.30, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1")
	assert.AssertEqualT(t, 3, len(e2), "e2")
	assert.AssertEqualT(t, 1, len(e3), "e3")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size")
	assert.AssertEqualT(t, 1, bk.BuyStopSize(), "buy stop size")
	assert.AssertEqualT(t, 0, bk.SellStopSize(), "sell stop size")

	e4, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id22", fixmodel.SideSell, 100, 1.50, aclock))
	assert.AssertEqualT(t, 4, len(e4), "e4")

	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size")
	assert.AssertEqualT(t, 0, bk.BuyStopSize(), "buy size")
	assert.AssertEqualT(t, 0, bk.SellStopSize(), "sell size")

	fixmodel.ContainsExec(t, e4, "cli2", "id22", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e4, "cli2", "id22", fixmodel.OrdStatusFilled, "filled", 100, 1.50, 1)
	fixmodel.ContainsExec(t, e4, "cli1", "id12", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e4, "cli1", "id12", fixmodel.OrdStatusFilled, "filled", 100, 1.50, 1)
}
