package orderbook

type OrderType rune

const (
	Market OrderType = '1'
	Limit            = '2'
	Stop             = '3'
)
