package orderbook

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
)

func Test_OrderBook_State_AuctionOpen_Closed(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenAuction, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeAuctionLimitOrder("cli1", "id1", SideBuy, 100, 1.01))
	e2, _ := bk.NewOrder(makeAuctionLimitOrder("cli2", "id2", SideSell, 101, 1.01))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, bk.BuyAuctionSize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellAuctionSize(), "sell size should be 1")
	containsExec(t, e1, "cli1", "id1", OrdStatusNew, "new", 0, 0)
	containsExec(t, e2, "cli2", "id2", OrdStatusNew, "new", 0, 0)

	e3, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 3, len(e3), "e3 1")
	containsExec(t, e3, "cli1", "id1", OrdStatusFilled, "fill", 100, 1.01)
	containsExec(t, e3, "cli2", "id2", OrdStatusPartiallyFilled, "part fill", 100, 1.01)
	containsExec(t, e3, "cli2", "id2", OrdStatusCanceled, "cancel", 0, 0)
}
