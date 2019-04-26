package orderbook

import (
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_AddBuySellOrder(t *testing.T) {
	ins := instrument.MakeInstrument("ABV", "ABV Investments")
	bk := MakeOrderBook(ins)
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	buy := MakeNewOrderEvent("id", 1.01, OrderTypeLimit, SideBuy, "")
	bk.NewOrder(buy)
	sell := MakeNewOrderEvent("id", 1.03, OrderTypeLimit, SideSell, "")
	bk.NewOrder(sell)

	assert.AssertEqualT(t, bk.BuySize(), 1, "buy size should be 1")
	assert.AssertEqualT(t, bk.SellSize(), 1, "sell size should be 1")
}
