package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/rdumont/assistdog"
	"orderbook/assert"
	"orderbook/instrument"
	"orderbook/orderbook"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var pending = fmt.Errorf("Pending")
var assit = assistdog.NewDefault()
var bk orderbook.OrderBook
var execs []orderbook.ExecutionReport
var orders []orderbook.OrderEvent
var loc, _ = time.LoadLocation("UTC")

const (
	tabEvent       string = "Event"
	tabClientID           = "ClientID"
	tabInstrument         = "Instrument"
	tabSide               = "Side"
	tabOrdType            = "OrdType"
	tabClOrdID            = "ClOrdID"
	tabOrigClOrdID        = "OrigClOrdID"
	tabPrice              = "Price"
	tabQty                = "Qty"
	tabExpireOn           = "ExpireOn"
	tabTimeInForce        = "TimeInForce"
	tabStatus             = "Status"
	tabExecType           = "ExecType"
	tabReason             = "Reason"
	tabExecID             = "ExecID"
	tabOrderID            = "OrderID"
	tabLastQty            = "LastQty"
	tabLastPrice          = "LastPrice"
	tabCumQty             = "CumQty"
)

func anOrderBookForInstrument(inst string) error {
	ins := instrument.MakeInstrument(inst, inst+"name")
	bk = orderbook.MakeOrderBook(ins)
	execs = []orderbook.ExecutionReport{}
	orders = []orderbook.OrderEvent{}
	return nil
}

func usersSendOrdersWith(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	for _, row := range slice {
		switch orderbook.EventTypeConv(row[tabEvent]) {
		case orderbook.EventTypeNewOrderSingle:
			order := makeOrder(row)
			exec, _ := bk.NewOrder(order)
			execs = append(execs, exec...)
			orders = append(orders, order)
		case orderbook.EventTypeCancel:
			order := makeCancelOrder(row)
			exec, _ := bk.CancelOrder(order)
			//fmt.Printf("Cancel Exec [%v]\n", exec)
			execs = append(execs, exec)
			orders = append(orders, order)
		}
	}
	return nil
}

func awaitExecutions(num int) error {
	if len(execs) == num {
		return nil
	}
	return fmt.Errorf("did not get %d execs, got %d instead", num, len(execs))
}

func executionsShouldBe(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	var expectedExecs []orderbook.ExecutionReport
	for _, row := range slice {
		exec := makeExec(row)
		expectedExecs = append(expectedExecs, exec)
	}
	for k, v := range expectedExecs {
		//fmt.Printf("Exp      Execs value[%s]\n", v)
		//fmt.Printf("Act k=%d Execs value[%s]\n", k, execs[k])
		if err := compareExec(v, execs[k]); err != nil {
			return err
		}

	}
	return nil
}

func FeatureContextLimitOrder(s *godog.Suite) {
	s.Step(`^An order book for instrument "([^"]*)"$`, anOrderBookForInstrument)
	s.Step(`^users send orders with:$`, usersSendOrdersWith)
	s.Step(`^await (\d+) executions$`, awaitExecutions)
	s.Step(`^executions should be:$`, executionsShouldBe)
}

func makeOrder(row map[string]string) orderbook.NewOrderSingle {
	price, _ := strconv.ParseFloat(row[tabPrice], 64)
	qty, _ := strconv.ParseInt(row[tabQty], 64, 64)
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	if row[tabOrdType] == "Limit" {
		return orderbook.MakeNewOrderLimit(
			row[tabInstrument],
			row[tabClientID],
			row[tabClOrdID],
			orderbook.SideConv(row[tabSide]),
			price,
			qty,
			orderbook.TimeInForceConv(row[tabTimeInForce]),
			dt,
			dt)
	}
	return nil

}

func makeCancelOrder(row map[string]string) orderbook.OrderCancelRequest {
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	return orderbook.MakeOrderCancelRequest(
		row[tabInstrument],
		row[tabClientID],
		row[tabClOrdID],
		orderbook.SideConv(row[tabSide]),
		row[tabOrigClOrdID],
		dt)
}

func makeExec(row map[string]string) orderbook.ExecutionReport {
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	qty, _ := strconv.ParseInt(row[tabQty], 64, 64)
	lastqty, _ := strconv.ParseInt(row[tabLastQty], 64, 64)
	cumqty, _ := strconv.ParseInt(row[tabCumQty], 64, 64)
	lastprice, _ := strconv.ParseFloat(row[tabLastPrice], 64)
	leavesQty := qty - cumqty
	return orderbook.MakeExecutionReport(
		orderbook.EventTypeConv(row[tabEvent]),
		row[tabInstrument],
		row[tabClientID],
		row[tabClOrdID],
		orderbook.SideConv(row[tabSide]),
		lastqty,
		lastprice,
		orderbook.ExecTypeConv(row[tabExecType]),
		leavesQty,
		cumqty,
		orderbook.OrdStatusConv(row[tabStatus]),
		row[tabOrderID],
		row[tabExecID],
		qty,
		dt)
}

func compareExec(exp orderbook.ExecutionReport, act orderbook.ExecutionReport) error {
	if act == nil {
		return fmt.Errorf("act nil")
	}
	if exp == nil {
		return fmt.Errorf("exp nil")
	}
	var errors strings.Builder
	assert.AssertEqualSB(exp.InstrumentID(), act.InstrumentID(), "InstrumentID", &errors)
	assert.AssertEqualSB(exp.ClientID(), act.ClientID(), "ClientID", &errors)
	assert.AssertEqualSB(exp.ClOrdID(), act.ClOrdID(), "ClOrdID", &errors)
	assert.AssertEqualSB(exp.Side(), act.Side(), "Side", &errors)
	assert.AssertEqualSB(exp.LastQty(), act.LastQty(), "LastQty", &errors)
	assert.AssertEqualSB(exp.LastPrice(), act.LastPrice(), "LastPrice", &errors)
	assert.AssertEqualSB(exp.ExecType(), act.ExecType(), "ExecType", &errors)
	assert.AssertEqualSB(exp.LeavesQty(), act.LeavesQty(), "LeavesQty", &errors)
	assert.AssertEqualSB(exp.CumQty(), act.CumQty(), "CumQty", &errors)
	assert.AssertEqualSB(exp.OrdStatus(), act.OrdStatus(), "OrdStatus", &errors)
	assert.AssertEqualSB(exp.OrderQty(), act.OrderQty(), "OrderQty", &errors)
	if !compareID(exp.OrderID(), act.OrderID()) {
		fmt.Fprintf(&errors, "%s", "orderid null")
	}
	if !compareID(exp.ExecID(), act.ExecID()) {
		fmt.Fprintf(&errors, "%s", "execid null")
	}
	if errors.Len() > 0 {
		return fmt.Errorf(errors.String())
	}
	return nil
}

func compareID(exp string, act string) bool {
	if exp == "Not Null" || exp == "Not Nil" {
		if reflect.TypeOf(act) == nil {
			return false
		}
		return true
	}
	return exp == act
}
