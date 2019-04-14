package orderbook

type Side int

const (
	Sell Side = -1
	Buy       = 1
)

type OrderType rune

const (
	Market OrderType = '1'
	Limit            = '2'
	Stop             = '3'
)
