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

func makeTime(hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(2019, 10, 11, hour, min, sec, 0, loc)
}

func makeLongTime(hour, min, sec int) time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Date(2019, 10, 21, hour, min, sec, 0, loc)
}

func makeLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64) NewOrderSingle {
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
		dt)
}

func makeLimitOrderWithClock(clientID string, clOrdID string, side Side, qty int64, price float64, clock clock.Clock) NewOrderSingle {
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

func makeLimitOrderTimeInForce(clientID string, clOrdID string, side Side, qty int64, price float64, timeInForce TimeInForce, expireOn time.Time) NewOrderSingle {
	dt := makeTime(11, 11, 1)
	return MakeNewOrderLimit(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		timeInForce,
		expireOn,
		dt)
}

func makeAuctionLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64) NewOrderSingle {
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
		dt)
}

func makeMarketOrder(clientID string, clOrdID string, side Side, qty int64) NewOrderSingle {
	dt := makeTime(11, 11, 1)
	return MakeNewOrderMarket(
		inst,
		clientID,
		clOrdID,
		side,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		dt)
}

func printExecs(execs []ExecutionReport) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
	fmt.Printf("i|clientid|clordid|side|lastprice|lastvol\n")
	for i, s := range execs {
		fmt.Printf("e%d|%s|%s|%s|%v|%v\n", i, s.ClientID(), s.ClOrdID(), SideToString(s.Side()), s.LastPrice(), s.LastQty())
	}
}

func containsExec(t *testing.T, execs []ExecutionReport, clientID string, clOrdID string, status OrdStatus, msg string, lastq int64, lastp float64) {
	var found = 0
	for _, v := range execs {
		if v.ClientID() == clientID &&
			v.ClOrdID() == clOrdID &&
			v.OrdStatus() == status &&
			v.LastPrice() == lastp &&
			v.LastQty() == lastq {
			found++
		}
	}
	if found == 0 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: not found %s %s:%s %v %d %f", assert.AssertionAt(file), line, msg, clientID, clOrdID, OrdStatusToString(status), lastq, lastp)
	}
}
