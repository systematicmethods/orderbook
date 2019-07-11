package orderbook

type OrderBookEventType int

type OrderBookState int

const (
	OrderBookStateTradingOpen OrderBookState = iota
	OrderBookStateTradingClosed
	OrderBookStateAuctionOpen
	OrderBookStateAuctionClosed
	OrderBookStateUnknown
)

const (
	OrderBookEventTypeOpenTrading OrderBookEventType = iota
	OrderBookEventTypeCloseTrading
	OrderBookEventTypeOpenAuction
	OrderBookEventTypeCloseAuction
	OrderBookEventTypeCloseAuctionAndOpenTrading
	OrderBookEventTypeUnknown
)

func OrderBookEventTypeAs(eventType OrderBookEventType) OrderBookState {
	switch eventType {
	case OrderBookEventTypeOpenTrading:
		return OrderBookStateTradingOpen
	case OrderBookEventTypeCloseTrading:
		return OrderBookStateTradingClosed
	case OrderBookEventTypeOpenAuction:
		return OrderBookStateAuctionOpen
	case OrderBookEventTypeCloseAuction:
		return OrderBookStateAuctionClosed
	}
	return OrderBookStateUnknown
}

func OrderBookStateConv(thetype string) OrderBookState {
	switch thetype {
	case "TradingOpen":
		return OrderBookStateTradingOpen
	case "TradingClosed":
		return OrderBookStateTradingClosed
	case "AuctionOpen":
		return OrderBookStateAuctionOpen
	case "AuctionClosed":
		return OrderBookStateAuctionClosed
	}
	return OrderBookStateUnknown
}

func OrderBookEventTypeConv(thetype string) OrderBookEventType {
	switch thetype {
	case "OpenTrading":
		return OrderBookEventTypeOpenTrading
	case "CloseTrading":
		return OrderBookEventTypeCloseTrading
	case "OpenAuction":
		return OrderBookEventTypeOpenAuction
	case "CloseAuction":
		return OrderBookEventTypeCloseAuction
	case "CloseAuctionAndOpenTrading":
		return OrderBookEventTypeCloseAuctionAndOpenTrading
	}
	return OrderBookEventTypeUnknown
}
