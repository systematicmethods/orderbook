package auction

import (
	"github.com/andres-erbsen/clock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"orderbook/fixmodel"
	"orderbook/instrument"
	"orderbook/orderstate"
	"orderbook/test"
	"orderbook/tradingevent"
)

func NewOrderBook(instrument instrument.Instrument, orderBookEvent tradingevent.OrderBookEventType, clock clock.Clock) *auction {
	b := auction{instrument: &instrument}
	b.orders.BuyOrders = orderstate.NewOrderList(orderstate.BuyPriceComparator)
	b.orders.SellOrders = orderstate.NewOrderList(orderstate.SellPriceComparator)
	b.orderBookState = tradingevent.OrderBookEventTypeAs(orderBookEvent)
	b.clock = clock
	// todo: yuk here just to import decimal
	priced, _ := decimal.NewFromString("1.23")
	priced.Add(priced)
	return &b
}

func NewMarketOrderForAuction(qty int64, price float64, side fixmodel.Side) *orderstate.OrderState {
	dt := test.NewTime(11, 11, 1)
	ordertype := fixmodel.OrderTypeMarket
	if price != 0 {
		ordertype = fixmodel.OrderTypeLimit
	}
	return orderstate.NewOrderState(
		"",
		"",
		"",
		side,
		price,
		qty,
		ordertype,
		fixmodel.TimeInForceGoodForAuction,
		dt,
		dt,
		dt,
		dt,
		"",
		uuid.New(),
		qty,
		0,
		fixmodel.OrdStatusNew,
	)
}

func NewAuctionLimitOrder(instrument string, clientID string, clOrdID string, side fixmodel.Side, qty int64, price float64, clock clock.Clock) *fixmodel.NewOrderSingle {
	dt := test.NewTime(11, 11, 1)
	return fixmodel.NewNewOrder(instrument,
		clientID,
		clOrdID,
		side,
		price,
		qty,
		fixmodel.TimeInForceGoodForAuction,
		dt,
		clock.Now(),
		fixmodel.OrderTypeLimit)
}
