package orderbookex

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_State_TradingOpen_Limit(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeLimitOrder("cli2", "id2", SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuySize(), 0, "buy size should be 0")
	assert.AssertEqualT(t, bk.SellSize(), 1, "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
}

func Test_OrderBook_State_TradingOpen_Market(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 101, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuySize(), 0, "buy size should be 0")
	assert.AssertEqualT(t, bk.SellSize(), 0, "sell size should be 0")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "new", 100, 1.01, 1)
}

func Test_OrderBook_State_TradingOpen_Auction(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id1", SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuyAuctionSize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellAuctionSize(), 1, "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0, 1)
}
