package orderbook

type EventType int

const (
	EventTypeNewOrderSingle EventType = iota
	EventTypeNewOrderAck
	EventTypeNewOrderRejected
	EventTypeCancel
	EventTypeCancelAck
	EventTypeCancelRejected
	EventTypeFill
	EventTypeDone
)
