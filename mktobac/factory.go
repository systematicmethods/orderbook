package mktobac

import (
	"github.com/andres-erbsen/clock"
	"orderbook/auction"
	"orderbook/instrument"
	"orderbook/orderbook"
	"orderbook/tradingevent"
)

func NewOrderBook(instrument instrument.Instrument, orderBookEvent tradingevent.OrderBookEventType, clock clock.Clock) *obac {
	b := obac{instrument: &instrument}
	b.ob = orderbook.NewOrderBook(instrument, orderBookEvent, clock)
	b.ac = auction.NewOrderBook(instrument, orderBookEvent, clock)
	b.clock = clock
	return &b
}
