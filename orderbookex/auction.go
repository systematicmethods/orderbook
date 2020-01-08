package orderbookex

import (
	"github.com/google/uuid"
	"time"
)

type OrderBookAuction interface {
	OpenAuction() error
	CloseAuction() (execs []ExecutionReport, clearingPrice float64, clearingVol int64, err error)
	BuyAuctionSize() int
	SellAuctionSize() int
	BuyAuctionOrders() []OrderState
	SellAuctionOrders() []OrderState

	auctionBookOrders() *buySellOrders
}

type buySellAuctionOrders struct {
	buyOrders  OrderList
	sellOrders OrderList
}

func (b *orderbook) OpenAuction() error {
	var err error
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeOpenAuction)
	return err
}

/*
	1	Find minimum price on buy side that match on the sell side
	2	Find maximum price on sell side that match on the buy side
	3	Find max volume that match between min and max price
	4	Match orders in price range and max volume on buy and sell side - use price time priority when orders are same price
	5	Calculate vwap buy orders using max volume
	6	Calculate vwap sell orders using max volume
	7	(vwap buy + vwap sell) /2
	8	Round to 2 decimals
	9	Fill orders to rounded auction price
	10	Cancel remaining orders
*/
func (b *orderbook) CloseAuction() (execs []ExecutionReport, clearingPrice float64, clearingVol int64, err error) {
	execs = []ExecutionReport{}
	state := newAuctionCloseCalculator()
	b.orderBookState, err = OrderBookStateChange(b.orderBookState, OrderBookEventTypeCloseAuction)
	if err == nil {
		var exs []ExecutionReport
		exs, err = state.fillAuctionAtClearingPrice(&b.auctionOrders)
		execs = append(execs, exs...)
		clearingPrice, _ = state.state().midclearingprice.Float64()
		clearingVol = state.state().clearingvol
		exs = cancelOrders(&b.auctionOrders)
		execs = append(execs, exs...)
	}
	return
}

func (b *orderbook) auctionBookOrders() *buySellOrders {
	return &b.auctionOrders
}

func (b *orderbook) BuyAuctionSize() int {
	return b.auctionOrders.buyOrders.Size()
}

func (b *orderbook) SellAuctionSize() int {
	return b.auctionOrders.sellOrders.Size()
}

func (b *orderbook) BuyAuctionOrders() []OrderState {
	return b.auctionOrders.buyOrders.Orders()
}

func (b *orderbook) SellAuctionOrders() []OrderState {
	return b.auctionOrders.sellOrders.Orders()
}

func makeMarketOrderForAuction(qty int64, price float64, side Side) OrderState {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 0, loc)
	ordertype := OrderTypeMarket
	if price != 0 {
		ordertype = OrderTypeLimit
	}
	return MakeOrderState(
		"",
		"",
		"",
		side,
		price,
		qty,
		ordertype,
		TimeInForceGoodForAuction,
		dt,
		dt,
		dt,
		dt,
		"",
		uuid.New(),
		qty,
		0,
		OrdStatusNew,
	)
}
