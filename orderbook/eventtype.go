package orderbook

type EventType int

const (
	NewOrderET EventType = iota
	NewOrderAckET
)
