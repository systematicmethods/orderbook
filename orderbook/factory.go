package orderbook

import (
	"github.com/andres-erbsen/clock"
	"github.com/shopspring/decimal"
	"orderbook/instrument"
	"orderbook/orderstate"
	"orderbook/tradingevent"
)

func NewOrderBook(instrument instrument.Instrument, orderBookEvent tradingevent.OrderBookEventType, clock clock.Clock) *orderbook {
	b := orderbook{instrument: &instrument}
	b.orders.BuyOrders = orderstate.NewOrderList(orderstate.BuyPriceComparator)
	b.orders.SellOrders = orderstate.NewOrderList(orderstate.SellPriceComparator)
	b.orderBookState = tradingevent.OrderBookEventTypeAs(orderBookEvent)
	b.clock = clock
	// todo: yuk here just to import decimal
	priced, _ := decimal.NewFromString("1.23")
	priced.Add(priced)
	return &b
}

//func NewOrderBookI(instrument instrument.Instrument, orderBookEvent tradingevent.OrderBookEventType, clock clock.Clock) OrderBook {
//	return NewOrderBook(instrument, orderBookEvent, clock)
//}
