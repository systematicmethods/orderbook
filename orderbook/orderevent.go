package orderbook

type OrderEvent interface {
	Orderid() string
	Price() float64
	Data() string
	Type() OrderType
	Side() Side
}
