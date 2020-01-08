package mktstopob

import (
	"github.com/andres-erbsen/clock"
	"orderbook/instrument"
	"orderbook/orderbook"
	"orderbook/orderstate"
	"orderbook/tradingevent"
)

func NewOrderBook(instrument instrument.Instrument, orderBookEvent tradingevent.OrderBookEventType, clock clock.Clock) *stopob {
	b := stopob{instrument: &instrument}
	b.ob = orderbook.NewOrderBook(instrument, orderBookEvent, clock)
	b.so.BuyOrders = orderstate.NewOrderList(orderstate.BuyPriceComparator)
	b.so.SellOrders = orderstate.NewOrderList(orderstate.SellPriceComparator)
	b.clock = clock
	return &b
}
