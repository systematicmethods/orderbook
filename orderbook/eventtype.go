package orderbook

type EventType int

const (
	EventTypeNewOrder EventType = iota
	EventTypeNewOrderAck
	EventTypeNewOrderRejected
	EventTypeCancel
	EventTypeCancelAck
	EventTypeCancelRejected
	EventTypeFill
	EventTypeDone
)
