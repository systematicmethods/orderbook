package orderbookex

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_IOC_Sell(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id2", SideSell, 101, 1.01, TimeInForceImmediateOrCancel, makeTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 4, len(e2), "e2 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled", 100, 1.01, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "part filled", 100, 1.01, 1)
	exec := containsExec(t, e2, "cli2", "id2", OrdStatusCanceled, "cancelled", 0, 0, 1)
	assert.AssertEqualT(t, exec.ExecType(), ExecTypeCanceled, "cancelled")
}

func Test_OrderBook_IOC_Buy(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 50, 1.01, aclock))
	e11, _ := bk.NewOrder(makeLimitOrder("cli1", "id2", SideSell, 50, 1.01, aclock))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id2", SideBuy, 101, 1.01, TimeInForceImmediateOrCancel, makeTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e11), "e1 1")
	assert.AssertEqualT(t, 6, len(e2), "e2 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled", 50, 1.01, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "part filled", 50, 1.01, 2)
	exec := containsExec(t, e2, "cli2", "id2", OrdStatusCanceled, "cancelled", 0, 0, 1)
	assert.AssertEqualT(t, exec.ExecType(), ExecTypeCanceled, "cancelled")
}
