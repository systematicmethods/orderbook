package orderbook

type OrderEvent interface {
	Orderid() string
	Price() float64
	Quantity() int64
	Data() string
	Type() OrderType
	Side() Side
	//UUID() uuid.UUID
}
