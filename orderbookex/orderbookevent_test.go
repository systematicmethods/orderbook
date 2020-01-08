package orderbookex

import (
	"orderbook/assert"
	"testing"
)

func Test_OrderBookChangeStateNoTradingToOpenTrading(t *testing.T) {
	var state, err = OrderBookStateChange(OrderBookStateOrderEntryClosed, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, OrderBookStateOrderEntryClosed, state, "open")
	assert.AssertNotNilT(t, err, "error null")
}

func Test_OrderBookChangeStateNoTradingToOrderEntryOpen(t *testing.T) {
	var state, err = OrderBookStateChange(OrderBookStateOrderEntryClosed, OrderBookEventTypeOpenOrderEntry)
	assert.AssertEqualT(t, OrderBookStateOrderEntryOpen, state, "open")
	assert.AssertNilT(t, err, "error null")
}

func Test_OrderBookChangeStateAuctionOpenToOpenTrading(t *testing.T) {
	var state, err = OrderBookStateChange(OrderBookStateAuctionOpen, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, OrderBookStateAuctionOpen, state, "open")
	assert.AssertNotNilT(t, err, "error not nil")
}
