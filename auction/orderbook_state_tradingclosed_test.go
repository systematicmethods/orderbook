package auction

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func Test_OrderBook_State_TradingClosed_Limit(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id2", fixmodel.SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "rej", 0, 0, 1)
}

func Test_OrderBook_State_TradingClosed_Market(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideSell, 101, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
}

func Test_OrderBook_State_TradingClosed_Auction(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(NewAuctionLimitOrder(inst, "cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(NewAuctionLimitOrder(inst, "cli2", "id2", fixmodel.SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "new", 0, 0, 1)
}

func Test_OrderBook_State_TradingClosed_Limit_FoK_IoC(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli1", "id1", fixmodel.SideBuy, 100, 1.01, fixmodel.TimeInForceFillOrKill, test.NewTime(11, 11, 1), aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli1", "id2", fixmodel.SideBuy, 100, 1.01, fixmodel.TimeInForceImmediateOrCancel, test.NewTime(11, 11, 1), aclock))
	e3, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id3", fixmodel.SideSell, 101, 1.01, fixmodel.TimeInForceFillOrKill, test.NewTime(11, 11, 1), aclock))
	e4, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id4", fixmodel.SideSell, 101, 1.01, fixmodel.TimeInForceImmediateOrCancel, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, len(e3), "e2 empty")
	assert.AssertEqualT(t, 1, len(e4), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id2", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id3", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
	fixmodel.ContainsExec(t, e4, "cli2", "id4", fixmodel.OrdStatusRejected, "rejected", 0, 0, 1)
}
