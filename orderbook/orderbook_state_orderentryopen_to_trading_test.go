package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
	"time"
)

func Test_OrderBook_State_OpenOrderEntry_Trading_Limit_SamePrice(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e11, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id2", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id2", fixmodel.SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e11), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 2, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)

	e3, _ := bk.OpenTrading()
	assert.AssertEqualT(t, 4, len(e3), "e3 empty")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e3, "cli1", "id1", fixmodel.OrdStatusFilled, "new", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli1", "id2", fixmodel.OrdStatusPartiallyFilled, "part", 1, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id2", fixmodel.OrdStatusFilled, "fill", 1, 1.01, 1)
}

func Test_OrderBook_State_OpenOrderEntry_Trading_Limit_SellLessThanBuy(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	aclock := test.NewMockClock(12, 34, 0)
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenOrderEntry, aclock)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	aclock.Add(time.Millisecond)
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id2", fixmodel.SideSell, 101, 1.00, aclock))
	aclock.Add(time.Millisecond)
	eb2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id3", fixmodel.SideBuy, 100, 1.01, aclock))
	aclock.Add(time.Millisecond)
	eb3, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id4", fixmodel.SideBuy, 100, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(eb2), "e1 1")
	assert.AssertEqualT(t, 1, len(eb3), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 3, bk.BuySize(), "buy size should be 3")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)

	e3, _ := bk.OpenTrading()

	//printExecs(e3)

	assert.AssertEqualT(t, 4, len(e3), "e3 empty")
	assert.AssertEqualT(t, 2, bk.BuySize(), "buy size should be 2")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e3, "cli1", "id1", fixmodel.OrdStatusFilled, "new", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli1", "id3", fixmodel.OrdStatusPartiallyFilled, "part", 1, 1.00, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id2", fixmodel.OrdStatusFilled, "fill", 1, 1.00, 1)
}
