package orderbook

import "errors"

var DuplicateOrder = errors.New("duplicate order")
var AuctionNotOpen = errors.New("invalid state: auction not open")
var AuctionOpen = errors.New("invalid state: auction open")
var TradingNotOpen = errors.New("invalid state: trading not open")
