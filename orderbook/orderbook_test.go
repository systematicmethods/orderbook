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
	bk := MakeOrderBook(ins)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeorder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeorder("cl12", "id2", SideSell, 101, 1.03))

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
	bk := MakeOrderBook(ins)

	e10, _ := bk.NewOrder(makeorder("cl12", "id2", SideBuy, 100, 1.01))
	e11, _ := bk.NewOrder(makeorder("cl12", "id3", SideBuy, 100, 1.01))
	e12, _ := bk.NewOrder(makeorder("cl12", "id4", SideBuy, 100, 1.01))
	e13, _ := bk.NewOrder(makeorder("cl12", "id5", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeorder("cli1", "id1", SideSell, 101, 1.00))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 5, len(e2), "e2 empty")
	assert.AssertEqualT(t, e2[0].OrdStatus(), OrdStatusNew, "new order")
	assert.AssertEqualT(t, e2[1].OrdStatus(), OrdStatusFilled, "filled")
	assert.AssertEqualT(t, e2[2].OrdStatus(), OrdStatusPartiallyFilled, "part filled")
	assert.AssertEqualT(t, 3, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 0, len(bk.SellOrders()), "sell orders")
}

func Test_OrderBook_MatchSellBuyOrder(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins)

	e10, _ := bk.NewOrder(makeorder("cl12", "id2", SideSell, 100, 1.00))
	e11, _ := bk.NewOrder(makeorder("cl12", "id3", SideSell, 100, 1.00))
	e12, _ := bk.NewOrder(makeorder("cl12", "id4", SideSell, 100, 1.00))
	e13, _ := bk.NewOrder(makeorder("cl12", "id5", SideSell, 100, 1.00))
	e2, _ := bk.NewOrder(makeorder("cli1", "id1", SideBuy, 101, 1.01))

	//printExecs(e2)

	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
	assert.AssertEqualT(t, 5, len(e2), "e2 empty")

	assert.AssertEqualT(t, e2[0].OrdStatus(), OrdStatusNew, "new order")
	assert.AssertEqualT(t, e2[1].OrdStatus(), OrdStatusPartiallyFilled, "part filled")
	assert.AssertEqualT(t, e2[2].OrdStatus(), OrdStatusFilled, "filled")
	assert.AssertEqualT(t, 0, len(bk.BuyOrders()), "buy orders")
	assert.AssertEqualT(t, 3, len(bk.SellOrders()), "sell orders")
}

func makeorder(clientID string, clOrdID string, side Side, qty int64, price float64) NewOrderSingle {
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

func printExecs(execs []ExecutionReport) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
}
