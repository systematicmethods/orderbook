package orderbookex

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_State_OpenOrderEntry_Limit(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeLimitOrder("cli2", "id2", SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
}

func Test_OrderBook_State_OpenOrderEntry_Market(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeMarketOrder("cli1", "id1", SideBuy, 100, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 101, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusRejected, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusRejected, "new", 0, 0, 1)
}

func Test_OrderBook_State_OpenOrderEntry_Auction(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, bk.BuyAuctionSize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellAuctionSize(), "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
}

func Test_OrderBook_State_OpenOrderEntry_Limit_FoK_IoC(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli1", "id1", SideBuy, 100, 1.01, TimeInForceFillOrKill, makeTime(11, 11, 1), aclock))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli1", "id2", SideBuy, 100, 1.01, TimeInForceImmediateOrCancel, makeTime(11, 11, 1), aclock))
	e3, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id3", SideSell, 101, 1.01, TimeInForceFillOrKill, makeTime(11, 11, 1), aclock))
	e4, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id4", SideSell, 101, 1.01, TimeInForceImmediateOrCancel, makeTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, len(e3), "e2 empty")
	assert.AssertEqualT(t, 1, len(e4), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusRejected, "rejected", 0, 0, 1)
	containsExec(t, e2, "cli1", "id2", OrdStatusRejected, "rejected", 0, 0, 1)
	containsExec(t, e3, "cli2", "id3", OrdStatusRejected, "rejected", 0, 0, 1)
	containsExec(t, e4, "cli2", "id4", OrdStatusRejected, "rejected", 0, 0, 1)
}
