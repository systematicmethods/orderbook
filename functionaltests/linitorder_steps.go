package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/rdumont/assistdog"
	"orderbook/instrument"
	"orderbook/orderbook"
	"reflect"
	"strconv"
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
	return nil
}

func usersSendOrdersWith(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	for _, row := range slice {
		switch orderbook.EventTypeConv(row[tabEvent]) {
		case orderbook.EventTypeNewOrderSingle:
			order := makeOrder(row)
			exec, _ := bk.NewOrder(order)
			execs = append(execs, exec)
			orders = append(orders, order)
		case orderbook.EventTypeCancel:
			order := makeCancelOrder(row)
			exec, _ := bk.CancelOrder(order)
			execs = append(execs, exec)
			orders = append(orders, order)
		}
	}
	return nil
}

func awaitExecutions(num int) error {
	if (bk.BuySize() + bk.SellSize()) == num {
		return nil
	}
	if len(execs) == num {
		return nil
	}
	return fmt.Errorf("did not get %d execs, got %d instead", num, (bk.BuySize() + bk.SellSize()))
}

func executionsShouldBe(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	var expectedExecs []orderbook.ExecutionReport
	for _, row := range slice {
		exec := makeExec(row)
		expectedExecs = append(expectedExecs, exec)
	}
	for k, v := range expectedExecs {
		fmt.Printf("Execs value[%s]\n", v)
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
		orderbook.ExecTypeConv(row[tabStatus]),
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
	if exp.InstrumentID() != act.InstrumentID() ||
		exp.ClientID() != act.ClientID() ||
		exp.ClOrdID() != act.ClOrdID() ||
		exp.Side() != act.Side() ||
		exp.LastQty() != act.LastQty() ||
		exp.LastPrice() != act.LastPrice() ||
		exp.ExecType() != act.ExecType() ||
		exp.LeavesQty() != act.LeavesQty() ||
		exp.CumQty() != act.CumQty() ||
		exp.OrdStatus() != act.OrdStatus() ||
		exp.OrderQty() != act.OrderQty() ||
		!compareID(exp.OrderID(), act.OrderID()) ||
		!compareID(exp.ExecID(), act.ExecID()) {
		return fmt.Errorf("Expect %v \nActual %v ", exp, act)
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
