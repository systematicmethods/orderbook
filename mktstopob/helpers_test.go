package mktstopob

import (
	"github.com/andres-erbsen/clock"
	"orderbook/fixmodel"
	"orderbook/test"
)

const inst = "ABV"

func makeStopOrder(clientID string, clOrdID string, side fixmodel.Side, qty int64, price float64, clock clock.Clock) *fixmodel.NewOrderSingle {
	dt := test.NewTime(11, 11, 1)
	return fixmodel.NewNewOrder(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		fixmodel.TimeInForceImmediateOrCancel,
		dt,
		clock.Now(),
		fixmodel.OrderTypeStop)
}
