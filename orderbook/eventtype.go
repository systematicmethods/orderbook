package orderbook

type EventType int

const (
	EventTypeNewOrder EventType = iota
	EventTypeNewOrderAck
)
