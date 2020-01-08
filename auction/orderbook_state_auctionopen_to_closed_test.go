package auction

import (
	clock "github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func Test_OrderBook_State_AuctionOpen_Closed(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenAuction, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(NewAuctionLimitOrder(inst, "cli1", "id1", fixmodel.SideBuy, 100, 1.01, aclock))
	e2, _ := bk.NewOrder(NewAuctionLimitOrder(inst, "cli2", "id2", fixmodel.SideSell, 101, 1.01, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 empty")
	assert.AssertEqualT(t, 1, len(e2), "e2 empty")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")
	fixmodel.ContainsExec(t, e1, "cli1", "id1", fixmodel.OrdStatusNew, "new", 0, 0, 1)
	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusNew, "new", 0, 0, 1)

	e3, clrPrice, clrVol, _ := bk.CloseAuction()
	assert.AssertEqualT(t, 3, len(e3), "e3 1")
	assert.AssertEqualT(t, 1.01, clrPrice, "clearing price")
	assert.AssertEqualT(t, int64(100), clrVol, "clearing vol")
	fixmodel.ContainsExec(t, e3, "cli1", "id1", fixmodel.OrdStatusFilled, "fill", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id2", fixmodel.OrdStatusPartiallyFilled, "part fill", 100, 1.01, 1)
	fixmodel.ContainsExec(t, e3, "cli2", "id2", fixmodel.OrdStatusCanceled, "cancel", 0, 0, 1)
	//printExecs(e3)

}
