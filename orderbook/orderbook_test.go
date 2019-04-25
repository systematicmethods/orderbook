package orderbook

import (
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_AddBuySellOrder(t *testing.T) {
	ins := instrument.MakeInstrument("ABV", "ABV Investments")
	bk := MakeOrderBook(ins)
	assert.AssertEqual(t, *bk.Instrument(), ins, "instrument same")

	buy := MakeNewOrderEvent("id", 1.01, Limit, Buy, "")
	bk.NewOrder(buy)
	sell := MakeNewOrderEvent("id", 1.03, Limit, Sell, "")
	bk.NewOrder(sell)

	assert.AssertEqual(t, bk.BuySize(), 1, "buy size should be 1")
	assert.AssertEqual(t, bk.SellSize(), 1, "sell size should be 1")
}
