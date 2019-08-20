package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_FillOrKill_Sell(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id2", SideSell, 101, 1.01, TimeInForceFillOrKill, makeTime(11, 11, 1)))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 2, len(e2), "e2 1")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e2, "cli2", "id2", OrdStatusRejected, "cancelled", 0, 0)
}

func Test_OrderBook_FillOrKill_Buy(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 100, 1.01))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id2", SideBuy, 101, 1.01, TimeInForceFillOrKill, makeTime(11, 11, 1)))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 2, len(e2), "e2 1")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 0")
	containsExec(t, e2, "cli2", "id2", OrdStatusRejected, "cancelled", 0, 0)
}
