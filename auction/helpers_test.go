package auction

import (
	"github.com/andres-erbsen/clock"
	"orderbook/fixmodel"
	"time"
)

const inst = "ABV"

var dt = time.Date(2019, 10, 11, 11, 11, 1, 0, time.UTC)

func newLimitOrder(clientID string, clOrdID string, side fixmodel.Side, qty int64, price float64, clock clock.Clock) *fixmodel.NewOrderSingle {
	return fixmodel.NewNewOrder(
		inst,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		fixmodel.TimeInForceGoodTillCancel,
		dt,
		clock.Now(),
		fixmodel.OrderTypeLimit)
}
