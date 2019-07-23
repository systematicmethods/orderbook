package orderbook

import (
	"orderbook/assert"
	"testing"
)

func Test_OrderBookChangeStateNoTradingToOpenTrading(t *testing.T) {
	var state, error = OrderBookStateChange(OrderBookStateNoTrading, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, OrderBookStateTradingOpen, state, "open")
	assert.AssertNilT(t, error, "error null")
}

func Test_OrderBookChangeStateAuctionOpenToOpenTrading(t *testing.T) {
	var state, error = OrderBookStateChange(OrderBookStateAuctionOpen, OrderBookEventTypeOpenTrading)
	assert.AssertEqualT(t, OrderBookStateAuctionOpen, state, "open")
	assert.AssertNotNilT(t, error, "error not nil")
}
