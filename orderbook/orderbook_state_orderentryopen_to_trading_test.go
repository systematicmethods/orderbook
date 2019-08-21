package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
	"time"
)

func Test_OrderBook_State_OpenOrderEntry_Trading_Limit_SamePrice(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e11, _ := bk.NewOrder(makeLimitOrder("cli1", "id2", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeLimitOrder("cli2", "id2", SideSell, 101, 1.01))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e11), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 2, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0)

	e3, _ := bk.OpenTrading()
	assert.AssertEqualT(t, 4, len(e3), "e3 empty")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e3, "cli1", "id1", OrdStatusFilled, "new", 100, 1.01)
	containsExec(t, e3, "cli1", "id2", OrdStatusPartiallyFilled, "part", 1, 1.01)
	containsExec(t, e3, "cli2", "id2", OrdStatusPartiallyFilled, "part", 100, 1.01)
	containsExec(t, e3, "cli2", "id2", OrdStatusFilled, "fill", 1, 1.01)
}

func Test_OrderBook_State_OpenOrderEntry_Trading_Limit_SellLessThanBuy(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	aclock := makeMockClock(12, 34, 0)
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenOrderEntry, aclock)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrderWithClock("cli1", "id1", SideBuy, 100, 1.01, aclock))
	aclock.Add(time.Millisecond)
	e2, _ := bk.NewOrder(makeLimitOrderWithClock("cli2", "id2", SideSell, 101, 1.00, aclock))
	aclock.Add(time.Millisecond)
	eb2, _ := bk.NewOrder(makeLimitOrderWithClock("cli1", "id3", SideBuy, 100, 1.01, aclock))
	aclock.Add(time.Millisecond)
	eb3, _ := bk.NewOrder(makeLimitOrderWithClock("cli1", "id4", SideBuy, 100, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(eb2), "e1 1")
	assert.AssertEqualT(t, 1, len(eb3), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 3, bk.BuySize(), "buy size should be 3")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0)

	e3, _ := bk.OpenTrading()

	printExecs(e3)

	assert.AssertEqualT(t, 4, len(e3), "e3 empty")
	assert.AssertEqualT(t, 2, bk.BuySize(), "buy size should be 2")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e3, "cli1", "id1", OrdStatusFilled, "new", 100, 1.01)
	containsExec(t, e3, "cli2", "id2", OrdStatusPartiallyFilled, "part", 100, 1.01)
	containsExec(t, e3, "cli1", "id3", OrdStatusPartiallyFilled, "part", 1, 1.00)
	containsExec(t, e3, "cli2", "id2", OrdStatusFilled, "fill", 1, 1.00)
}
