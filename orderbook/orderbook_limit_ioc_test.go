package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func Test_OrderBook_IOC_Sell(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id2", fixmodel.SideSell, 101, 1.01, fixmodel.TimeInForceImmediateOrCancel, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 4, len(e2), "e2 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "filled", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part filled", 100, 1.01, 1)
	exec := fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusCanceled, "cancelled", 0, 0, 1)
	assert.AssertEqualT(t, exec.ExecType(), fixmodel.ExecTypeCanceled, "cancelled")
}

func Test_OrderBook_IOC_Buy(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideSell, 50, 1.01, aclock))
	e11, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id2", fixmodel.SideSell, 50, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id2", fixmodel.SideBuy, 101, 1.01, fixmodel.TimeInForceImmediateOrCancel, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e11), "e1 1")
	assert.AssertEqualT(t, 6, len(e2), "e2 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "filled", 50, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part filled", 50, 1.01, 2)
	exec := fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusCanceled, "cancelled", 0, 0, 1)
	assert.AssertEqualT(t, exec.ExecType(), fixmodel.ExecTypeCanceled, "cancelled")
}
