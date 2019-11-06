package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_RejectBuySellMarketOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeMarketOrder("cli1", "id1", SideBuy, 100, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cl12", "id2", SideSell, 101, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuySize(), 0, "buy size should be 0")
	assert.AssertEqualT(t, bk.SellSize(), 0, "sell size should be 0")
	assert.AssertEqualT(t, e1[0].OrdStatus(), OrdStatusRejected, "should be rejected")
	assert.AssertEqualT(t, e2[0].OrdStatus(), OrdStatusRejected, "should be rejected")
}

func Test_OrderBook_MatchBuyLimitSellMarket(t *testing.T) {
	loglines()
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 100, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")

	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusPartiallyFilled, "part filled", 100, 1.0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusFilled, "filled", 100, 1.0, 1)

	//printExecs(e2)
}

func Test_OrderBook_MatchSellLimitBuyMarket(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideBuy, 100, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 1, len(bk.SellOrders()), "sell orders")

	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusPartiallyFilled, "part filled", 100, 1.0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusFilled, "filled", 100, 1.0, 1)
}

func Test_OrderBook_MatchBuyLimitSellMarketLeaves(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 105, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")

	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled", 101, 1.0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "part filled", 101, 1.0, 1)
}

func Test_OrderBook_MatchSellLimitBuyMarketLeaves(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := makeMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideBuy, 105, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0, 1)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled", 101, 1.0, 1)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "part filled", 101, 1.0, 1)
}
