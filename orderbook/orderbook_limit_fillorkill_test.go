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

func Test_OrderBook_FillOrKill_Sell(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id2", fixmodel.SideSell, 101, 1.01, fixmodel.TimeInForceFillOrKill, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 2, len(e2), "e2 1")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "cancelled", 0, 0, 1)
}

func Test_OrderBook_FillOrKill_Buy(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	aclock := test.NewMockClock(12, 34, 0)
	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideSell, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id2", fixmodel.SideBuy, 101, 1.01, fixmodel.TimeInForceFillOrKill, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 2, len(e2), "e2 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "cancelled", 0, 0, 1)
}
