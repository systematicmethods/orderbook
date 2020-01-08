package mktobac

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/auction"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func Test_OrderBook_State_CloseOrderEntry_Limit(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id2", fixmodel.SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e2 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "rej", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "rej", 0, 0, 1)
}

func Test_OrderBook_State_CloseOrderEntry_Market(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli1", "id1", fixmodel.SideBuy, 100, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideSell, 101, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "new", 0, 0, 1)
}

func Test_OrderBook_State_CloseOrderEntry_Auction(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeCloseOrderEntry, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(auction.NewAuctionLimitOrder(inst, "cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(auction.NewAuctionLimitOrder(inst, "cli2", "id2", fixmodel.SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 1")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusRejected, "rej", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusRejected, "rej", 0, 0, 1)
}
