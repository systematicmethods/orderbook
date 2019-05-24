package orderbook

import "time"

type OrderEvent interface {
	InstrumentID() string
	ClientID() string
	ClOrdID() string
	Side() Side
	OrderID() string
	TransactTime() time.Time
}
