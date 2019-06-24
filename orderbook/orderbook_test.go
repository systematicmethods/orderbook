package orderbook

import (
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

	buy := makeorder("id1", SideBuy, 100, 1.01)
	e1, _ := bk.NewOrder(buy)
	sell := makeorder("id2", SideSell, 101, 1.03)
	e2, _ := bk.NewOrder(sell)

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, bk.BuySize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellSize(), 1, "sell size should be 1")
	assert.AssertEqualT(t, e1[0].ClOrdID(), "id1", "same clord")
	assert.AssertEqualT(t, e2[0].ClOrdID(), "id2", "same clord")
	assert.AssertEqualT(t, e1[0].InstrumentID(), inst, "same instrument")
	assert.AssertEqualT(t, e2[0].InstrumentID(), inst, "same instrument")
}

func makeorder(clOrdID string, side Side, qty int64, price float64) NewOrderSingle {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	return MakeNewOrderLimit(
		inst,
		"clientID",
		clOrdID,
		side,
		price,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		dt)
}
