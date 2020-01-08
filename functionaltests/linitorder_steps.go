package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"github.com/rdumont/assistdog"
	"orderbook/assert"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/orderbook"
	"orderbook/orderstate"
	"orderbook/tradingevent"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var pending = fmt.Errorf("Pending")
var assit = assistdog.NewDefault()
var bk orderbook.OrderBook
var execs []*fixmodel.ExecutionReport
var orders []fixmodel.OrderEvent
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
	ins := instrument.NewInstrument(inst, inst+"name")
	bk = orderbook.NewOrderBook(ins, tradingevent.OrderBookEventTypeOpenTrading, clock.NewMock())
	execs = []*fixmodel.ExecutionReport{}
	orders = []fixmodel.OrderEvent{}
	return nil
}

func usersSendOrdersWith(table *gherkin.DataTable) error {
	execs = []*fixmodel.ExecutionReport{}
	slice, _ := assit.ParseSlice(table)
	for _, row := range slice {
		switch fixmodel.EventTypeConv(row[tabEvent]) {
		case fixmodel.EventTypeNewOrderSingle:
			order := newOrder(row)
			executions, _ := bk.NewOrder(order)
			execs = append(execs, executions...)
			orders = append(orders, order)
		case fixmodel.EventTypeCancel:
			order := newCancelOrder(row)
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

func containsExec(ex []*fixmodel.ExecutionReport, ac *fixmodel.ExecutionReport, msg string) error {
	if ac == nil {
		printExecs(msg, ex)
		return fmt.Errorf("Act nil %s", msg)

	}
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
	var expectedExecs []*fixmodel.ExecutionReport
	for _, row := range slice {
		exec := newExec(row)
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
	var expectedBuyState []*orderstate.OrderState
	var expectedSellState []*orderstate.OrderState
	for _, row := range slice {
		order := newState(row)
		if order.Side() == fixmodel.SideBuy {
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

func printState(msg string, orders []*orderstate.OrderState) {
	for _, v := range orders {
		fmt.Printf("%s: %v\n", msg, v)
	}
}

func printExecs(msg string, execs []*fixmodel.ExecutionReport) {
	for _, v := range execs {
		fmt.Printf("%s: %v\n", msg, v)
	}
}

func newOrder(row map[string]string) *fixmodel.NewOrderSingle {
	price, _ := strconv.ParseFloat(row[tabPrice], 64)
	qty, _ := strconv.ParseInt(row[tabQty], 10, 64)
	//fmt.Printf("\n qty in test %d %f row %v\n", qty, price, row)
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	if row[tabOrdType] == "Limit" {
		return fixmodel.NewNewOrder(
			row[tabInstrument],
			row[tabClientID],
			row[tabClOrdID],
			fixmodel.SideConv(row[tabSide]),
			price,
			qty,
			fixmodel.TimeInForceConv(row[tabTimeInForce]),
			dt,
			dt,
			fixmodel.OrderTypeLimit)
	} else if row[tabOrdType] == "Market" {
		return fixmodel.NewNewOrder(
			row[tabInstrument],
			row[tabClientID],
			row[tabClOrdID],
			fixmodel.SideConv(row[tabSide]),
			0,
			qty,
			fixmodel.TimeInForceConv(row[tabTimeInForce]),
			dt,
			dt,
			fixmodel.OrderTypeMarket)
	}
	return nil

}

func newCancelOrder(row map[string]string) *fixmodel.OrderCancelRequest {
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	return fixmodel.NewOrderCancelRequest(
		row[tabInstrument],
		row[tabClientID],
		row[tabClOrdID],
		fixmodel.SideConv(row[tabSide]),
		row[tabOrigClOrdID],
		dt)
}

func newExec(row map[string]string) *fixmodel.ExecutionReport {
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	qty, _ := strconv.ParseInt(row[tabQty], 10, 64)
	lastqty, _ := strconv.ParseInt(row[tabLastQty], 10, 64)
	cumqty, _ := strconv.ParseInt(row[tabCumQty], 10, 64)
	lastprice, _ := strconv.ParseFloat(row[tabLastPrice], 64)
	leavesQty := qty - cumqty
	return fixmodel.NewExecutionReport(
		fixmodel.EventTypeConv(row[tabEvent]),
		row[tabInstrument],
		row[tabClientID],
		row[tabClOrdID],
		fixmodel.SideConv(row[tabSide]),
		lastqty,
		lastprice,
		fixmodel.ExecTypeConv(row[tabExecType]),
		leavesQty,
		cumqty,
		fixmodel.OrdStatusConv(row[tabStatus]),
		row[tabOrigClOrdID],
		row[tabOrderID],
		row[tabExecID],
		qty,
		dt,
		fixmodel.ExecRestatementReasonNone)
}

func newState(row map[string]string) *orderstate.OrderState {
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	qty, _ := strconv.ParseInt(row[tabQty], 10, 64)
	cumqty, _ := strconv.ParseInt(row[tabCumQty], 10, 64)
	price, _ := strconv.ParseFloat(row[tabPrice], 64)
	leavesQty := qty - cumqty
	return orderstate.NewOrderState(
		row[tabInstrument],
		row[tabClientID],
		row[tabClOrdID],
		fixmodel.SideConv(row[tabSide]),
		price,
		qty,
		fixmodel.OrderTypeConv(row[tabOrdType]),
		fixmodel.TimeInForceConv(row[tabTimeInForce]),
		dt,
		dt,
		dt,
		dt,
		row[tabOrderID],
		uuid.New(),
		leavesQty,
		cumqty,
		fixmodel.OrdStatusConv(row[tabStatus]),
	)
}

func compareState(exp *orderstate.OrderState, act *orderstate.OrderState) error {
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
	assert.AssertEqualSB(fixmodel.OrdStatusToString(exp.OrdStatus()), fixmodel.OrdStatusToString(act.OrdStatus()), "OrdStatus", &errors)
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
func compareExec(exp *fixmodel.ExecutionReport, act *fixmodel.ExecutionReport) error {
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
	assert.AssertEqualSB(fixmodel.OrdStatusToString(exp.OrdStatus()), fixmodel.OrdStatusToString(act.OrdStatus()), "OrdStatus", &errors)
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
