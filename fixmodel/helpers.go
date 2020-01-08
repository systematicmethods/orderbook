package fixmodel

import (
	"encoding/csv"
	"fmt"
	"github.com/andres-erbsen/clock"
	"github.com/gocarina/gocsv"
	"io"
	"orderbook/assert"
	"orderbook/test"
	"runtime"
	"strings"
	"testing"
	"time"
)

const inst = "ABV"

func PrintExecs(execs []*ExecutionReport) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
	fmt.Printf("id|clientid|clordid|side|lastprice|lastqty|status\n")
	for i, s := range execs {
		fmt.Printf("e%d|%s|%s|%s|%v|%v|%v\n", i, s.ClientID(), s.ClOrdID(), s.Side(), s.LastPrice(), s.LastQty(), OrdStatusToString(s.OrdStatus()))
	}
}

func ContainsExec(t *testing.T, execs []*ExecutionReport, clientID string, clOrdID string, status OrdStatus, msg string, lastq int64, lastp float64, count int) *ExecutionReport {
	var found = 0
	var exec *ExecutionReport
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
		PrintExecs(execs)
	} else if found != count {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: found (%d) too many %s %s:%s %v %d %f", assert.AssertionAt(file), line, found, msg, clientID, clOrdID, OrdStatusToString(status), lastq, lastp)
		PrintExecs(execs)
	}
	return exec
}

func NewLimitOrder(clientID string, clOrdID string, side Side, qty int64, price float64, clock clock.Clock) *NewOrderSingle {
	dt := test.NewTime(11, 11, 1)
	return NewNewOrder(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		clock.Now(),
		OrderTypeLimit)
}

func NewLimitOrderTimeInForce(clientID string, clOrdID string, side Side, qty int64, price float64, timeInForce TimeInForce, expireOn time.Time, clock clock.Clock) *NewOrderSingle {
	return NewNewOrder(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		timeInForce,
		expireOn,
		clock.Now(),
		OrderTypeLimit)
}

func NewMarketOrder(clientID string, clOrdID string, side Side, qty int64, clock clock.Clock) *NewOrderSingle {
	dt := test.NewTime(11, 11, 1)
	return NewNewOrder(
		inst,
		clientID,
		clOrdID,
		side,
		0,
		qty,
		TimeInForceGoodTillCancel,
		dt,
		clock.Now(),
		OrderTypeMarket)
}

type execCSV struct {
	Id        string    `csv:"id"`
	Clientid  string    `csv:"clientid"`
	Clordid   string    `csv:"clordid"`
	Side      Side      `csv:"side"`
	Lastprice float64   `csv:"lastprice"`
	Lastqty   int64     `csv:"lastqty"`
	Status    OrdStatus `csv:"status"`
}

func (ex *Side) MarshalCSV() (string, error) {
	return ex.String(), nil
}

func (ex *Side) UnmarshalCSV(field string) (err error) {
	*ex = SideConv(field)
	return nil
}

func (ex *OrdStatus) MarshalCSV() (string, error) {
	return ex.String(), nil
}

func (ex *OrdStatus) UnmarshalCSV(field string) (err error) {
	*ex = OrdStatusConv(field)
	return nil
}

func LoadExecCSV(csvstr string) []*execCSV {
	execs := []*execCSV{}

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		return r
	})

	if gocsv.Unmarshal(strings.NewReader(csvstr), &execs) == nil {
		return execs
	}
	return execs
}

func ContainsExecCSV(t *testing.T, execs []*ExecutionReport, ex *execCSV, msg string) {
	var found = 0
	for _, v := range execs {
		if v.ClientID() == ex.Clientid &&
			v.ClOrdID() == ex.Clordid &&
			v.OrdStatus() == ex.Status &&
			v.LastPrice() == ex.Lastprice &&
			v.LastQty() == ex.Lastqty {
			found++
		}
	}
	if found == 0 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("\n%s:%d: not found %s id:'%v' %s:%s stat:%v qty:%d price:%f",
			assert.AssertionAt(file), line, msg, ex.Id, ex.Clientid, ex.Clordid, ex.Status, ex.Lastqty, ex.Lastprice)
	}
}
