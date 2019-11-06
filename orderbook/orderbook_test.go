package orderbook

import (
	"fmt"
	"github.com/andres-erbsen/clock"
	"orderbook/assert"
	"runtime"
	"testing"
	"time"
)

const inst = "ABV"

func makeMockClock(hour, min, sec int) *clock.Mock {
	aclock := clock.NewMock()
	aclock.Set(makeTime(hour, min, sec))
	return aclock
}

func makeMockClockFromTime(dt time.Time) *clock.Mock {
	aclock := clock.NewMock()
	aclock.Set(dt)
	return aclock
}
func makeTime(hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(2019, 10, 11, hour, min, sec, 0, loc)
}

func makeLongTime(hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(2019, 10, 21, hour, min, sec, 0, loc)
}

func makeDateTime(year, month, day, hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
}

func makeLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64, clock clock.Clock) NewOrderSingle {
	dt := makeTime(11, 11, 1)
	return MakeNewOrderLimit(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		clock.Now())
}

func makeLimitOrderTimeInForce(clientID string, clOrdID string, side Side, qty int64, price float64, timeInForce TimeInForce, expireOn time.Time, clock clock.Clock) NewOrderSingle {
	return MakeNewOrderLimit(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		timeInForce,
		expireOn,
		clock.Now())
}

func makeAuctionLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64, clock clock.Clock) NewOrderSingle {
	dt := makeTime(11, 11, 1)
	return MakeNewOrderLimit(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		TimeInForceGoodForAuction,
		dt,
		clock.Now())
}

func makeMarketOrder(clientID string, clOrdID string, side Side, qty int64, clock clock.Clock) NewOrderSingle {
	dt := makeTime(11, 11, 1)
	return MakeNewOrderMarket(
		inst,
		clientID,
		clOrdID,
		side,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		clock.Now())
}

func printExecs(execs []ExecutionReport) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
	fmt.Printf("id|clientid|clordid|side|lastprice|lastqty|status\n")
	for i, s := range execs {
		fmt.Printf("e%d|%s|%s|%s|%v|%v|%v\n", i, s.ClientID(), s.ClOrdID(), s.Side(), s.LastPrice(), s.LastQty(), OrdStatusToString(s.OrdStatus()))
	}
}

func printOrders(bk *buySellOrders) {
	fmt.Printf("clientid|clordid|side|gty|price \n")
	for iter := bk.buyOrders.iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		fmt.Printf("%s|%s|%v|%v|%v \n", order.ClientID(), order.ClOrdID(), order.Side(), order.OrderQty(), order.Price())
	}
	for iter := bk.sellOrders.iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		fmt.Printf("%s|%s|%v|%v|%v \n", order.ClientID(), order.ClOrdID(), order.Side(), order.OrderQty(), order.Price())
	}
}

func printExecsAndOrders(execs []ExecutionReport, bk *buySellOrders, buyorders []OrderState, sellorders []OrderState) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
	fmt.Printf("id|clientid|clordid|side|lastprice|lastqty|status|price|qty|ordstatus\n")
	for i, ex := range execs {
		var order OrderState
		if ex.Side() == SideBuy {
			for _, anorder := range buyorders {
				if anorder.ClOrdID() == ex.ClOrdID() {
					order = anorder
				}
			}
		} else {
			for _, anorder := range sellorders {
				if anorder.ClOrdID() == ex.ClOrdID() {
					order = anorder
				}
			}
		}
		if order != nil {
			fmt.Printf("e%d|%s|%s|%s|%v|%v|%v|%v|%v|%v\n", i, ex.ClientID(), ex.ClOrdID(), ex.Side(), ex.LastPrice(), ex.LastQty(), ex.OrdStatus(), order.Price(), order.OrderQty(), order.OrdStatus())
		} else {
			fmt.Printf("e%d|%s|%s|%s|%v|%v|%v|\n", i, ex.ClientID(), ex.ClOrdID(), ex.Side(), ex.LastPrice(), ex.LastQty(), ex.OrdStatus())
		}
	}
}

func containsExec(t *testing.T, execs []ExecutionReport, clientID string, clOrdID string, status OrdStatus, msg string, lastq int64, lastp float64, count int) ExecutionReport {
	var found = 0
	var exec ExecutionReport
	for _, v := range execs {
		if v.ClientID() == clientID &&
			v.ClOrdID() == clOrdID &&
			v.OrdStatus() == status &&
			v.LastPrice() == lastp &&
			v.LastQty() == lastq {
			found++
			exec = v
		}
	}
	if found == 0 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: not found %s %s:%s %v %d %f", assert.AssertionAt(file), line, msg, clientID, clOrdID, OrdStatusToString(status), lastq, lastp)
		printExecs(execs)
	}
	if found != count {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: found (%d) too many %s %s:%s %v %d %f", assert.AssertionAt(file), line, found, msg, clientID, clOrdID, OrdStatusToString(status), lastq, lastp)
		printExecs(execs)
	}
	return exec
}
