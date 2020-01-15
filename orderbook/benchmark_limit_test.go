package orderbook

import (
	"fmt"
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/test"
	"orderbook/tradingevent"
	"testing"
)

func BenchmarkXxx(b *testing.B) {
	ins := instrument.NewInstrument(inst, "ABV Investments")
	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	aclock := test.NewMockClock(12, 34, 0)
	for i := 0; i < b.N; i++ {
		id1 := fmt.Sprintf("id1%d", i)
		id2 := fmt.Sprintf("id2%d", i)
		bk.NewOrder(fixmodel.NewLimitOrder("cli1", id1, fixmodel.SideBuy, 100, 1.01, aclock))
		bk.NewOrder(fixmodel.NewLimitOrder("cl12", id2, fixmodel.SideSell, 101, 1.03, aclock))
	}
	println("")
	assert.AssertEqual(bk.BuySize(), 100, "buy size")
	assert.AssertEqual(bk.SellSize(), 100, "sell size")
}

func BenchmarkHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

//
//func Test_OrderBook_Perf_MatchBuySellOrder(t *testing.T) {
//	ins := instrument.NewInstrument(inst, "ABV Investments")
//	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
//	aclock := test.NewMockClock(12, 34, 0)
//
//	e10, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id2", fixmodel.SideBuy, 100, 1.01, aclock))
//	e11, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id3", fixmodel.SideBuy, 100, 1.01, aclock))
//	e12, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id4", fixmodel.SideBuy, 100, 1.01, aclock))
//	e13, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id5", fixmodel.SideBuy, 100, 1.01, aclock))
//	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id1", fixmodel.SideSell, 101, 1.00, aclock))
//
//	//printExecs(e2)
//
//	assert.AssertEqualT(t, 1, len(e10), "e10 1")
//	assert.AssertEqualT(t, 1, len(e11), "e11 1")
//	assert.AssertEqualT(t, 1, len(e12), "e12 1")
//	assert.AssertEqualT(t, 1, len(e13), "e13 1")
//	assert.AssertEqualT(t, 5, len(e2), "e2 5")
//	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
//	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusPartiallyFilled, "partially filled order", 100, 1.01, 1)
//	fixmodel.ContainsExec(t, e2, "cli1", "id1", fixmodel.OrdStatusFilled, "filled order", 1, 1.01, 1)
//	fixmodel.ContainsExec(t, e2, "cli2", "id2", fixmodel.OrdStatusFilled, "filled bk order", 100, 1.01, 1)
//	fixmodel.ContainsExec(t, e2, "cli2", "id3", fixmodel.OrdStatusPartiallyFilled, "partially filled bk order", 1, 1.01, 1)
//
//}
//
//func Test_OrderBook_Perf_MatchSellBuyOrder(t *testing.T) {
//	ins := instrument.NewInstrument(inst, "ABV Investments")
//	bk := NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
//	aclock := test.NewMockClock(12, 34, 0)
//
//	e10, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id2", fixmodel.SideSell, 100, 1.00, aclock))
//	e11, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id3", fixmodel.SideSell, 100, 1.00, aclock))
//	e12, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id4", fixmodel.SideSell, 100, 1.00, aclock))
//	e13, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli1", "id5", fixmodel.SideSell, 100, 1.00, aclock))
//	e2, _ := bk.NewOrder(fixmodel.NewLimitOrder("cli2", "id1", fixmodel.SideBuy, 101, 1.01, aclock))
//
//	assert.AssertEqualT(t, 1, len(e10), "e10 empty")
//	assert.AssertEqualT(t, 1, len(e11), "e11 empty")
//	assert.AssertEqualT(t, 1, len(e12), "e12 empty")
//	assert.AssertEqualT(t, 1, len(e13), "e13 empty")
//	assert.AssertEqualT(t, 5, len(e2), "e2 empty")
//
//	//printExecs(e2)
//
//	fixmodel.ContainsExec(t, e2, "cli2", "id1", fixmodel.OrdStatusNew, "new order", 0, 0, 1)
//	fixmodel.ContainsExec(t, e2, "cli2", "id1", fixmodel.OrdStatusPartiallyFilled, "partially filled order", 100, 1.00, 1)
//	fixmodel.ContainsExec(t, e2, "cli2", "id1", fixmodel.OrdStatusFilled, "filled order", 1, 1.00, 1)
//	fixmodel.ContainsExec(t, e2, "cli1", "id2", fixmodel.OrdStatusFilled, "filled bk order", 100, 1.00, 1)
//	fixmodel.ContainsExec(t, e2, "cli1", "id3", fixmodel.OrdStatusPartiallyFilled, "partially filled bk order", 1, 1.00, 1)
//}
