package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"orderbook/instrument"
	"orderbook/orderbook"
)

var pending = fmt.Errorf("Pending")

var bk orderbook.OrderBook

func anOrderBookForInstrument(inst string) error {
	ins := instrument.MakeInstrument(inst, inst+"name")
	bk = orderbook.MakeOrderBook(ins)
	return nil
}

func usersSendOrdersWith(table *gherkin.DataTable) error {
	for _, row := range table.Rows[1:] {
		order := orderbook.MakeNewOrderEvent(row.Cells[1].Value, 1, orderbook.Limit, orderbook.Buy, "")
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
	for _, row := range table.Rows[1:] {
		order := orderbook.MakeNewOrderEvent(row.Cells[1].Value, 1, orderbook.Limit, orderbook.Buy, "")
		bk.NewOrder(order)
	}
	return godog.ErrPending
}

func FeatureContextLimitOrder(s *godog.Suite) {
	s.Step(`^An order book for instrument "([^"]*)"$`, anOrderBookForInstrument)
	s.Step(`^users send orders with:$`, usersSendOrdersWith)
	s.Step(`^await (\d+) executions$`, awaitExecutions)
	s.Step(`^executions should be:$`, executionsShouldBe)
}

func makeOrderFromRow(row *gherkin.TableRow) orderbook.OrderEvent {
	return orderbook.MakeNewOrderEvent(row.Cells[1].Value, 1, orderbook.Limit, orderbook.Buy, "")
}
