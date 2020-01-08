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

func Test_OrderBook_RejectBuySellMarketOrder(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli1", "id1", fixmodel.SideBuy, 100, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cl12", "id2", fixmodel.SideSell, 101, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuySize(), 0, "buy size should be 0")
	assert.AssertEqualT(t, bk.SellSize(), 0, "sell size should be 0")
	assert.AssertEqualT(t, e1[0].OrdStatus(), fixmodel.OrdStatusRejected, "should be rejected")
	assert.AssertEqualT(t, e2[0].OrdStatus(), fixmodel.OrdStatusRejected, "should be rejected")
}

func Test_OrderBook_MatchBuyLimitSellMarket(t *testing.T) {
	test.Loglines()
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideSell, 100, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")

	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusPartiallyFilled, "part filled", 100, 1.0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusFilled, "filled", 100, 1.0, 1)

	//printExecs(e2)
}

func Test_OrderBook_MatchSellLimitBuyMarket(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideSell, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideBuy, 100, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 1, len(bk.SellOrders()), "sell orders")

	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusPartiallyFilled, "part filled", 100, 1.0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusFilled, "filled", 100, 1.0, 1)
}

func Test_OrderBook_MatchBuyLimitSellMarketLeaves(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideBuy, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideSell, 105, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")

	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "filled", 101, 1.0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part filled", 101, 1.0, 1)
}

func Test_OrderBook_MatchSellLimitBuyMarketLeaves(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideSell, 101, 1.00, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewMarketOrder("cli2", "id2", fixmodel.SideBuy, 105, aclock))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "filled", 101, 1.0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part filled", 101, 1.0, 1)
}
