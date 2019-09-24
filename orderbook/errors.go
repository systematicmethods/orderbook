package orderbook

type Error string

func (e Error) Error() string { return string(e) }

const DuplicateOrder = Error("duplicate order")
const AuctionNotOpen = Error("invalid state: auction not open")
const AuctionOpen = Error("invalid state: auction open")
const TradingNotOpen = Error("invalid state: trading not open")
const OrderBookEmpty = Error("orderbook empty")
