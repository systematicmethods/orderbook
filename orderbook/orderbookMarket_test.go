package orderbook

import (
	"orderbook/assert"
	"orderbook/instrument"
	"runtime"
	"testing"
)

func Test_OrderBook_RejectBuySellMarketOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeMarketOrder("cli1", "id1", SideBuy, 100))
	e2, _ := bk.NewOrder(makeMarketOrder("cl12", "id2", SideSell, 101))

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
	bk := MakeOrderBook(ins)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 101, 1.00))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 100))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")

	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusPartiallyFilled, "part filled", 100, 1.0)
	containsExec(t, e2, "cli2", "id2", OrdStatusFilled, "filled", 100, 1.0)

	//printExecs(e2)
}

func containsExec(t *testing.T, execs []ExecutionReport, clientID string, clOrdID string, status OrdStatus, msg string, lastq int64, lastp float64) {
	var found = 0
	for _, v := range execs {
		if v.ClientID() == clientID &&
			v.ClOrdID() == clOrdID &&
			v.OrdStatus() == status &&
			v.LastPrice() == lastp &&
			v.LastQty() == lastq {
			found++
		}
	}
	if found == 0 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: not found %s %s:%s %v %d %f", assert.AssertionAt(file), line, msg, clientID, clOrdID, status, lastq, lastp)
	}
}

func Test_OrderBook_MatchSellLimitBuyMarket(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 101, 1.00))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideBuy, 100))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 1, len(bk.SellOrders()), "sell orders")

	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusPartiallyFilled, "part filled", 100, 1.0)
	containsExec(t, e2, "cli2", "id2", OrdStatusFilled, "filled", 100, 1.0)
}

func Test_OrderBook_MatchBuyLimitSellMarketLeaves(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideBuy, 101, 1.00))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideSell, 105))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")

	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled", 101, 1.0)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "part filled", 101, 1.0)
}

func Test_OrderBook_MatchSellLimitBuyMarketLeaves(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins)

	e1, _ := bk.NewOrder(makeLimitOrder("cli1", "id1", SideSell, 101, 1.00))
	e2, _ := bk.NewOrder(makeMarketOrder("cli2", "id2", SideBuy, 105))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 3, len(e2), "e2 empty")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new order", 0, 0)
	containsExec(t, e2, "cli1", "id1", OrdStatusFilled, "filled", 101, 1.0)
	containsExec(t, e2, "cli2", "id2", OrdStatusPartiallyFilled, "part filled", 101, 1.0)
}
