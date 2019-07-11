package orderbook

import (
	"fmt"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
	"time"
)

const inst = "ABV"

func Test_OrderBook_AddBuySellOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeLimitOrder("cl12", "id2", SideSell, 101, 1.03))

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
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading)

	e10, _ := bk.NewOrder(makeLimitOrder("cli2", "id2", SideBuy, 100, 1.01))
	e11, _ := bk.NewOrder(makeLimitOrder("cli2", "id3", SideBuy, 100, 1.01))
	e12, _ := bk.NewOrder(makeLimitOrder("cli2", "id4", SideBuy, 100, 1.01))
	e13, _ := bk.NewOrder(makeLimitOrder("cli2", "id5", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 101, 1.00))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e10), "e10 1")
	assert.AssertEqualT(t, 1, len(e11), "e11 1")
	assert.AssertEqualT(t, 1, len(e12), "e12 1")
	assert.AssertEqualT(t, 1, len(e13), "e13 1")
	assert.AssertEqualT(t, 5, len(e2), "e2 5")
	containsExec(t, e2, "cli1", "id1", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusPartiallyFilled, "partially filled order", 100, 1.01)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled order", 1, 1.01)
	containsExec(t, e2, "cli2", "id2", OrdStatusFilled, "filled bk order", 100, 1.01)
	containsExec(t, e2, "cli2", "id3", OrdStatusPartiallyFilled, "partially filled bk order", 1, 1.01)

}

func Test_OrderBook_MatchSellBuyOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading)

	e10, _ := bk.NewOrder(makeLimitOrder("cli1", "id2", SideSell, 100, 1.00))
	e11, _ := bk.NewOrder(makeLimitOrder("cli1", "id3", SideSell, 100, 1.00))
	e12, _ := bk.NewOrder(makeLimitOrder("cli1", "id4", SideSell, 100, 1.00))
	e13, _ := bk.NewOrder(makeLimitOrder("cli1", "id5", SideSell, 100, 1.00))
	e2, _ := bk.NewOrder(makeLimitOrder("cli2", "id1", SideBuy, 101, 1.01))

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 5, len(e2), "e2 empty")

	//printExecs(e2)

	containsExec(t, e2, "cli2", "id1", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli2", "id1", OrdStatusPartiallyFilled, "partially filled order", 100, 1.00)
	containsExec(t, e2, "cli2", "id1", OrdStatusFilled, "filled order", 1, 1.00)
	containsExec(t, e2, "cli1", "id2", OrdStatusFilled, "filled bk order", 100, 1.00)
	containsExec(t, e2, "cli1", "id3", OrdStatusPartiallyFilled, "partially filled bk order", 1, 1.00)
}

func Test_OrderBook_TradingClosed(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeCloseTrading)

	e10, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 100, 1.00))
	assert.AssertEqualT(t, 1, len(e10), "e10")
	containsExec(t, e10, "cli1", "id1", OrdStatusRejected, "rejected", 0, 0)

	e11, _ := bk.NewOrder(makeLimitOrder("cli2", "id1", SideSell, 100, 1.00))
	assert.AssertEqualT(t, 1, len(e11), "e11")
	containsExec(t, e11, "cli2", "id1", OrdStatusRejected, "rejected", 0, 0)
}

func makeLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64) NewOrderSingle {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	return MakeNewOrderLimit(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		dt)
}

func makeMarketOrder(clientID string, clOrdID string, side Side, qty int64) NewOrderSingle {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	return MakeNewOrderMarket(
		inst,
		clientID,
		clOrdID,
		side,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		dt)
}

func printExecs(execs []ExecutionReport) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
}
