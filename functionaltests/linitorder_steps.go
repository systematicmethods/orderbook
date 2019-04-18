package functionaltests

import (
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

var pending = fmt.Errorf("Pending")

func anOrderBookForInstrument(arg1 string) error {
	return pending
}

func usersSendOrdersWith(arg1 *gherkin.DataTable) error {
	return godog.ErrUndefined
}

func awaitExecutions(arg1 int) error {
	return godog.ErrPending
}

func executionsShouldBe(arg1 *gherkin.DataTable) error {
	return godog.ErrPending
}

func FeatureContextLimitOrder(s *godog.Suite) {
	s.Step(`^An order book for instrument "([^"]*)"$`, anOrderBookForInstrument)
	s.Step(`^users send orders with:$`, usersSendOrdersWith)
	s.Step(`^await (\d+) executions$`, awaitExecutions)
	s.Step(`^executions should be:$`, executionsShouldBe)
}
