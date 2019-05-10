package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/rdumont/assistdog"
	"orderbook/assert"
	"orderbook/instrument"
	"orderbook/orderbook"
	"strconv"
)

var pending = fmt.Errorf("Pending")
var assit = assistdog.NewDefault()
var bk orderbook.OrderBook

func anOrderBookForInstrument(inst string) error {
	ins := instrument.MakeInstrument(inst, inst+"name")
	bk = orderbook.MakeOrderBook(ins)
	return nil
}

func usersSendOrdersWith(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	for _, row := range slice {
		order := makeOrder(row)
		bk.NewOrder(order)
	}
	return nil
}

func awaitExecutions(num int) error {
	if (bk.BuySize() + bk.SellSize()) == num {
		return nil
	}
	return fmt.Errorf("did not get %d execs, got %d instead", num, (bk.BuySize() + bk.SellSize()))
}

func executionsShouldBe(table *gherkin.DataTable) error {
	slice, _ := assit.ParseSlice(table)
	for _, row := range slice {
		var other orderbook.Order
		order := makeOrder(row)
		if order.Side() == orderbook.SideSell {
			other = bk.SellOrders()[0]
		} else {
			other = bk.BuyOrders()[0]

		}
		if err := assert.AssertEqual(other.Orderid(), order.OrderID(), "orderid should be the same"); err != nil {
			return err
		}
		if err := assert.AssertEqual(other.Price(), order.Price(), "price should be the same"); err != nil {
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

func makeOrder(row map[string]string) orderbook.OrderEvent {
	price, _ := strconv.ParseFloat(row["Price"], 64)
	return orderbook.MakeNewOrderEvent(row["ClOrdID"],
		price,
		orderbook.OrderTypeConv(row["OrdType"]),
		orderbook.SideConv(row["Side"]),
		"")
}
