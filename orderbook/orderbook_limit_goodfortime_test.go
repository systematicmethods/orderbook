package orderbook

import (
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
	"time"
)

func Test_OrderBook_AddBuySellOrderGoodForTime(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	e1, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli1", "id1", fixmodel.SideBuy, 100, 1.01, fixmodel.TimeInForceGoodForTime, test.NewTime(11, 11, 1), aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id2", fixmodel.SideSell, 101, 1.03, fixmodel.TimeInForceGoodForTime, test.NewTime(11, 11, 1), aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")

	clock := test.NewMockClock(11, 11, 0)
	e3, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 0, len(e3), "e3 empty")

	clock.Add(time.Second)
	e4, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 2, len(e4), "e4 2")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
}

func Test_OrderBook_AddBuySellOrderGoodForLongTime(t *testing.T) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	assert.AssertEqualT(t, *bk.Instrument(), ins, "instrument same")
	aclock := test.NewMockClock(12, 34, 0)

	expireon := test.NewDateTime(2019, 10, 21, 11, 11, 11)
	e1, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli1", "id1", fixmodel.SideBuy, 100, 1.01, fixmodel.TimeInForceGoodForTime, expireon, aclock))
	e2, _ := bk.NewOrder(fixmodel.NewLimitOrderTimeInForce("cli2", "id2", fixmodel.SideSell, 101, 1.03, fixmodel.TimeInForceGoodForTime, expireon, aclock))

	assert.AssertEqualT(t, 1, len(e1), "e1 1")
	assert.AssertEqualT(t, 1, len(e2), "e2 1")
	assert.AssertEqualT(t, 1, bk.BuySize(), "buy size should be 1")
	assert.AssertEqualT(t, 1, bk.SellSize(), "sell size should be 1")

	obtime := test.NewDateTime(2019, 10, 21, 11, 11, 1)
	clock := test.NewMockClockFromTime(obtime)
	e3, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 0, len(e3), "e3 empty")

	clock.Add(time.Second)
	e4, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 0, len(e4), "e4 2")

	clock.Add(time.Minute)
	e5, _ := bk.Tick(clock.Now())
	assert.AssertEqualT(t, 2, len(e5), "e5 2")
	assert.AssertEqualT(t, 0, bk.BuySize(), "buy size should be 0")
	assert.AssertEqualT(t, 0, bk.SellSize(), "sell size should be 0")
}
