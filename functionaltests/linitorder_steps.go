package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/google/uuid"
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
	tabEvent        string = "Event"
	tabClientID            = "ClientID"
	tabInstrument          = "Instrument"
	tabSide                = "Side"
	tabOrdType             = "OrdType"
	tabClOrdID             = "ClOrdID"
	tabOrigClOrdID         = "OrigClOrdID"
	tabPrice               = "Price"
	tabQty                 = "Qty"
	tabLeavesQty           = "LeavesQty"
	tabExpireOn            = "ExpireOn"
	tabTimeInForce         = "TimeInForce"
	tabStatus              = "Status"
	tabExecType            = "ExecType"
	tabReason              = "Reason"
	tabExecID              = "ExecID"
	tabOrderID             = "OrderID"
	tabLastQty             = "LastQty"
	tabLastPrice           = "LastPrice"
	tabCumQty              = "CumQty"
	tabTransactTime        = "TransactTime"
	tabCreatedOn           = "CreatedOn"
	tabUpdatedOn           = "UpdatedOn"
	tabTimestamp           = "Timestamp"
)

func FeatureContextLimitOrder(s *godog.Suite) {
	s.Step(`^An order book for instrument "([^"]*)"$`, anOrderBookForInstrument)
	s.Step(`^users send orders with:$`, usersSendOrdersWith)
	s.Step(`^await (\d+) executions$`, awaitExecutions)
	s.Step(`^executions should be:$`, executionsShouldBe)
	s.Step(`^order state should be:$`, orderStateShouldBe)
}

func anOrderBookForInstrument(inst string) error {
	ins := instrument.MakeInstrument(inst, inst+"name")
	bk = orderbook.MakeOrderBook(ins, orderbook.OrderBookEventTypeOpenTrading)
	execs = []orderbook.ExecutionReport{}
	orders = []orderbook.OrderEvent{}
	return nil
}

func usersSendOrdersWith(table *gherkin.DataTable) error {
	execs = []orderbook.ExecutionReport{}
	slice, _ := assit.ParseSlice(table)
	for _, row := range slice {
		switch orderbook.EventTypeConv(row[tabEvent]) {
		case orderbook.EventTypeNewOrderSingle:
			order := makeOrder(row)
			executions, _ := bk.NewOrder(order)
			execs = append(execs, executions...)
			orders = append(orders, order)
		case orderbook.EventTypeCancel:
			order := makeCancelOrder(row)
			executions, _ := bk.CancelOrder(order)
			//fmt.Printf("Cancel Exec [%v]\n", exec)
			execs = append(execs, executions)
			orders = append(orders, order)
		}
	}
	return nil
}

func awaitExecutions(num int) error {
	if len(execs) == num {
		return nil
	}
	for exec := range execs {
		fmt.Printf("awaitExecutions %v\n", execs[exec])
	}
	return fmt.Errorf("did not get %d execs, got %d", num, len(execs))
}

func containsExec(ex []orderbook.ExecutionReport, ac orderbook.ExecutionReport, msg string) error {
	var found = 0
	for _, v := range ex {
		if err := compareExec(v, ac); err == nil {
			found++
		}
	}
	if found == 0 {
		printExecs(msg, ex)
		return fmt.Errorf("%s %v", msg, ac)
	}
	return nil
}

func executionsShouldBe(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	var expectedExecs []orderbook.ExecutionReport
	for _, row := range slice {
		exec := makeExec(row)
		expectedExecs = append(expectedExecs, exec)
	}
	for _, ac := range execs {
		if err := containsExec(expectedExecs, ac, "exec not found"); err != nil {
			return err
		}
	}
	return nil
}

func orderStateShouldBe(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	var expectedBuyState []orderbook.OrderState
	var expectedSellState []orderbook.OrderState
	for _, row := range slice {
		order := makeState(row)
		if order.Side() == orderbook.SideBuy {
			expectedBuyState = append(expectedBuyState, order)
		} else {
			expectedSellState = append(expectedSellState, order)
		}
	}

	var errors strings.Builder
	buyOrders := bk.BuyOrders()
	if !assert.AssertEqualSB(len(expectedBuyState), len(buyOrders), "buy order state len different", &errors) {
		printState("expectedBuyState", expectedBuyState)
		printState("buyOrders", buyOrders)
	}
	sellOrders := bk.SellOrders()
	if !assert.AssertEqualSB(len(expectedSellState), len(sellOrders), "sell order state len different", &errors) {
		printState("expectedSellState", expectedSellState)
		printState("sellOrders", sellOrders)
	}
	if errors.Len() > 0 {
		return fmt.Errorf(errors.String())
	}

	for k, v := range expectedBuyState {
		if err := compareState(v, buyOrders[k]); err != nil {
			return err
		}
	}

	for k, v := range expectedSellState {
		if err := compareState(v, sellOrders[k]); err != nil {
			return err
		}
	}
	return nil
}

func printState(msg string, orders []orderbook.OrderState) {
	for _, v := range orders {
		fmt.Printf("%s: %v\n", msg, v)
	}
}

func printExecs(msg string, execs []orderbook.ExecutionReport) {
	for _, v := range execs {
		fmt.Printf("%s: %v\n", msg, v)
	}
}

func makeOrder(row map[string]string) orderbook.NewOrderSingle {
	price, _ := strconv.ParseFloat(row[tabPrice], 64)
	qty, _ := strconv.ParseInt(row[tabQty], 10, 64)
	//fmt.Printf("\n qty in test %d %f row %v\n", qty, price, row)
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
	} else if row[tabOrdType] == "Market" {
		return orderbook.MakeNewOrderMarket(
			row[tabInstrument],
			row[tabClientID],
			row[tabClOrdID],
			orderbook.SideConv(row[tabSide]),
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
	qty, _ := strconv.ParseInt(row[tabQty], 10, 64)
	lastqty, _ := strconv.ParseInt(row[tabLastQty], 10, 64)
	cumqty, _ := strconv.ParseInt(row[tabCumQty], 10, 64)
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
		row[tabOrigClOrdID],
		row[tabOrderID],
		row[tabExecID],
		qty,
		dt)
}

func makeState(row map[string]string) orderbook.OrderState {
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	qty, _ := strconv.ParseInt(row[tabQty], 10, 64)
	cumqty, _ := strconv.ParseInt(row[tabCumQty], 10, 64)
	price, _ := strconv.ParseFloat(row[tabPrice], 64)
	leavesQty := qty - cumqty
	return orderbook.MakeOrderState(
		row[tabInstrument],
		row[tabClientID],
		row[tabClOrdID],
		orderbook.SideConv(row[tabSide]),
		price,
		qty,
		orderbook.OrderTypeConv(row[tabOrdType]),
		orderbook.TimeInForceConv(row[tabTimeInForce]),
		dt,
		dt,
		dt,
		dt,
		row[tabOrderID],
		uuid.New(),
		leavesQty,
		cumqty,
		orderbook.OrdStatusConv(row[tabStatus]),
	)
}

func compareState(exp orderbook.OrderState, act orderbook.OrderState) error {
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
	assert.AssertEqualSB(exp.Price(), act.Price(), "Price", &errors)
	assert.AssertEqualSB(exp.LeavesQty(), act.LeavesQty(), "LeavesQty", &errors)
	assert.AssertEqualSB(exp.CumQty(), act.CumQty(), "CumQty", &errors)
	assert.AssertEqualSB(orderbook.OrdStatusToString(exp.OrdStatus()), orderbook.OrdStatusToString(act.OrdStatus()), "OrdStatus", &errors)
	assert.AssertEqualSB(exp.OrderQty(), act.OrderQty(), "OrderQty", &errors)
	if !compareID(exp.OrderID(), act.OrderID()) {
		fmt.Fprintf(&errors, "%s", "orderid null")
	}
	if !compareID(exp.OrderID(), act.OrderID()) {
		fmt.Fprintf(&errors, "%s", "OrderID null")
	}
	if errors.Len() > 0 {
		return fmt.Errorf(errors.String())
	}
	return nil
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
	assert.AssertEqualSB(orderbook.OrdStatusToString(exp.OrdStatus()), orderbook.OrdStatusToString(act.OrdStatus()), "OrdStatus", &errors)
	assert.AssertEqualSB(exp.OrigClOrdID(), act.OrigClOrdID(), "OrigClOrdID", &errors)
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
