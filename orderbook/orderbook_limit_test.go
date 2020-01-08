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

func Test_OrderBook_AddBuySellOrder(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cl12", "id2", fixmodel.SideSell, 101, 1.03, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuySize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellSize(), 1, "sell size should be 1")
	assert.AssertEqualT(t, e1[0].ClOrdID(), "id1", "same clord")
	assert.AssertEqualT(t, e2[0].ClOrdID(), "id2", "same clord")
	assert.AssertEqualT(t, e1[0].InstrumentID(), inst, "same instrument")
	assert.AssertEqualT(t, e2[0].InstrumentID(), inst, "same instrument")
}

func Test_OrderBook_MatchBuySellOrder(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)

	e10, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id2", fixmodel.SideBuy, 100, 1.01, aclock))
	e11, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id3", fixmodel.SideBuy, 100, 1.01, aclock))
	e12, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id4", fixmodel.SideBuy, 100, 1.01, aclock))
	e13, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id5", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideSell, 101, 1.00, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e10), "e10 1")
	assert.AssertEqualT(t, 1, len(e11), "e11 1")
	assert.AssertEqualT(t, 1, len(e12), "e12 1")
	assert.AssertEqualT(t, 1, len(e13), "e13 1")
	assert.AssertEqualT(t, 5, len(e2), "e2 5")
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusPartiallyFilled, "partially filled order", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "filled order", 1, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusFilled, "filled bk order", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id3", fixmodel.OrdStatusPartiallyFilled, "partially filled bk order", 1, 1.01, 1)

}

func Test_OrderBook_MatchSellBuyOrder(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)

	e10, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id2", fixmodel.SideSell, 100, 1.00, aclock))
	e11, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id3", fixmodel.SideSell, 100, 1.00, aclock))
	e12, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id4", fixmodel.SideSell, 100, 1.00, aclock))
	e13, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id5", fixmodel.SideSell, 100, 1.00, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id1", fixmodel.SideBuy, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 5, len(e2), "e2 empty")

	//printExecs(e2)

	fixmodel.ContainsExec(t, e2, "cli2", "id1", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id1", fixmodel.OrdStatusPartiallyFilled, "partially filled order", 100, 1.00, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id1", fixmodel.OrdStatusFilled, "filled order", 1, 1.00, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id2", fixmodel.OrdStatusFilled, "filled bk order", 100, 1.00, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id3", fixmodel.OrdStatusPartiallyFilled, "partially filled bk order", 1, 1.00, 1)
}
