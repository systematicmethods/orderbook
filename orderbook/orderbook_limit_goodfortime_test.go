package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/instrument"
	"testing"
	"time"
)

func Test_OrderBook_AddBuySellOrderGoodForTime(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli1", "id1", SideBuy, 100, 1.01, TimeInForceGoodForTime, makeTime(11, 11, 1)))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id2", SideSell, 101, 1.03, TimeInForceGoodForTime, makeTime(11, 11, 1)))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")

	clock := makeMockClock(11, 11, 0)
	e3, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 0, len(e3), "e3 empty")

	clock.Add(time.Second)
	e4, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 2, len(e4), "e4 2")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
}

func Test_OrderBook_AddBuySellOrderGoodForLongTime(t *testing.T) {
	ins := instrument.MakeInstrument(inst, "ABV Investments")
	bk := MakeOrderBook(ins, OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")

	e1, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli1", "id1", SideBuy, 100, 1.01, TimeInForceGoodForTime, makeLongTime(11, 11, 1)))
	e2, _ := bk.NewOrder(makeLimitOrderTimeInForce("cli2", "id2", SideSell, 101, 1.03, TimeInForceGoodForTime, makeLongTime(11, 11, 1)))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")

	clock := makeMockClock(11, 11, 0)
	e3, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 0, len(e3), "e3 empty")

	clock.Add(time.Second)
	e4, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 0, len(e4), "e4 2")

	clock.Add(time.Until(makeLongTime(11, 11, 1)))
	e5, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 2, len(e5), "e5 2")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
}
