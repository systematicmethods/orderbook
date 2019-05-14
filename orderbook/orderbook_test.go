package orderbook

import (
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
	"time"
)

func Test_OrderBook_AddBuySellOrder(t *testing.T) {
	ins := instrument.MakeInstrument("ABV", "ABV Investments")
	bk := MakeOrderBook(ins)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	buy := makeorder("id", SideBuy, 100, 1.01)
	bk.NewOrder(buy)
	sell := makeorder("id", SideSell, 101, 1.03)
	bk.NewOrder(sell)

	assert.AssertEqualT(t, bk.BuySize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellSize(), 1, "sell size should be 1")
}

func makeorder(clOrdID string, side Side, qty int64, price float64) NewOrderSingle {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	return MakeNewOrderLimit(
		"instrumentID",
		"clientID",
		"clOrdID",
		side,
		price,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		dt)
}

//func Test_OrderBook_CancelOrder(t *testing.T) {
//	ins := instrument.MakeInstrument("ABV", "ABV Investments")
//	bk := MakeOrderBook(ins)
//	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
//
//	buy := MakeNewOrderLimit("id", 1.01, OrderTypeLimit, SideBuy, "")
//	bk.NewOrder2(buy)
//
//	assert.AssertEqualT(t, bk.BuySize(), 1, "buy size should be 1")
//	t.Errorf("Pending")
//}
