package orderbook

import "fmt"

type OrderBookEventType int

type OrderBookState int

const (
	OrderBookStateTradingOpen OrderBookState = iota
	OrderBookStateTradingClosed
	OrderBookStateAuctionOpen
	OrderBookStateAuctionClosed
	OrderBookStateOrderEntryOpen
	OrderBookStateOrderEntryClosed
	OrderBookStateUnknown
)

const (
	OrderBookEventTypeOpenTrading OrderBookEventType = iota
	OrderBookEventTypeCloseTrading
	OrderBookEventTypeOpenAuction
	OrderBookEventTypeCloseAuction
	OrderBookEventTypeOpenOrderEntry
	OrderBookEventTypeCloseOrderEntry
	OrderBookEventTypeUnknown
)

var OrderBookEventTypeStateErrorFormat = "invalid event: %s in state %s"

func OrderBookStateChange(orderBookState OrderBookState, eventType OrderBookEventType) (OrderBookState, error) {
	switch eventType {
	case OrderBookEventTypeOpenTrading:
		if orderBookState == OrderBookStateOrderEntryOpen || orderBookState == OrderBookStateAuctionClosed {
			return OrderBookStateTradingOpen, nil
		}
		return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))
	case OrderBookEventTypeCloseTrading:
		if orderBookState == OrderBookStateTradingOpen {
			return OrderBookStateTradingClosed, nil
		}
		return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))

	case OrderBookEventTypeOpenAuction:
		if orderBookState == OrderBookStateOrderEntryOpen || orderBookState == OrderBookStateTradingClosed {
			return OrderBookStateAuctionOpen, nil
		}
		return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))
	case OrderBookEventTypeCloseAuction:
		if orderBookState == OrderBookStateAuctionOpen {
			return OrderBookStateAuctionClosed, nil
		}
		return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))

	case OrderBookEventTypeCloseOrderEntry:
		if orderBookState == OrderBookStateTradingClosed || orderBookState == OrderBookStateAuctionClosed {
			return OrderBookStateOrderEntryClosed, nil
		}
		return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))
	case OrderBookEventTypeOpenOrderEntry:
		if orderBookState == OrderBookStateOrderEntryClosed || orderBookState == OrderBookStateAuctionClosed {
			return OrderBookStateOrderEntryOpen, nil
		}
		return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))
	}
	return orderBookState, fmt.Errorf(OrderBookEventTypeStateErrorFormat, OrderBookEventTypeToString(eventType), OrderBookStateToString(orderBookState))
}

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
	case OrderBookEventTypeOpenOrderEntry:
		return OrderBookStateOrderEntryOpen
	case OrderBookEventTypeCloseOrderEntry:
		return OrderBookStateOrderEntryClosed
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
	case "OrderEntryClosed":
		return OrderBookStateOrderEntryClosed
	case "OrderEntryOpen":
		return OrderBookStateOrderEntryOpen
	}
	return OrderBookStateUnknown
}

func OrderBookStateToString(thetype OrderBookState) string {
	switch thetype {
	case OrderBookStateTradingOpen:
		return "TradingOpen"
	case OrderBookStateTradingClosed:
		return "TradingClosed"
	case OrderBookStateAuctionOpen:
		return "AuctionOpen"
	case OrderBookStateAuctionClosed:
		return "AuctionClosed"
	case OrderBookStateOrderEntryClosed:
		return "OrderEntryClosed"
	case OrderBookStateOrderEntryOpen:
		return "OrderEntryOpen"
	}
	return "OrderBookStateUnknown"
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
	case "OpenOrderEntry":
		return OrderBookEventTypeOpenOrderEntry
	case "CloseOrderEntry":
		return OrderBookEventTypeCloseOrderEntry
	}
	return OrderBookEventTypeUnknown
}

func OrderBookEventTypeToString(thetype OrderBookEventType) string {
	switch thetype {
	case OrderBookEventTypeOpenTrading:
		return "OpenTrading"
	case OrderBookEventTypeCloseTrading:
		return "CloseTrading"
	case OrderBookEventTypeOpenAuction:
		return "OpenAuction"
	case OrderBookEventTypeCloseAuction:
		return "CloseAuction"
	case OrderBookEventTypeOpenOrderEntry:
		return "OpenOrderEntry"
	case OrderBookEventTypeCloseOrderEntry:
		return "CloseOrderEntry"
	}
	return "OrderBookEventTypeUnknown"
}
