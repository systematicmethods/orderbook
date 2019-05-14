package orderbook

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
)

type OrderOfList int

const (
	LowToHigh OrderOfList = -1 // Sell orders increase in price
	HighToLow             = 1  // Buy orders decrease in price
)

type OrderList interface {
	Add(order Order) error
	Top() Order
	RemoveByID(orderid string) bool
	FindByID(orderid string) Order
	FindByPrice(price float64) []Order
	Orders() []Order
	Size() int
}

type orderlist struct {
	orderedlist *treeset.Set
	ordermap    map[string]Order
}

func NewOrderListStruct(sort OrderOfList) *orderlist {
	p := orderlist{}
	if sort == LowToHigh {
		p.orderedlist = treeset.NewWith(sellPriceComparator)
	} else if sort == HighToLow {
		p.orderedlist = treeset.NewWith(buyPriceComparator)
	}
	p.ordermap = make(map[string]Order)
	return &p
}

func NewOrderList(sort OrderOfList) OrderList {
	return NewOrderListStruct(sort)
}

func (p *orderlist) Add(order Order) error {
	if ord := p.ordermap[order.OrderID()]; ord != nil {
		return DuplicateOrder
	}
	p.orderedlist.Add(order)
	p.ordermap[order.OrderID()] = order
	return nil
}

func (p *orderlist) Size() int {
	return p.orderedlist.Size()
}

func (p *orderlist) RemoveByID(orderid string) bool {
	if ord := p.FindByID(orderid); ord != nil {
		p.orderedlist.Remove(ord)
		delete(p.ordermap, orderid)
		return true
	}
	return false
}

func (p *orderlist) Top() Order {
	var iter = p.orderedlist.Iterator()
	iter.Next()
	return iter.Value().(Order)
}

func (p *orderlist) Orders() []Order {
	var orders []Order
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(Order)
		orders = append(orders, order)
	}
	return orders
}

func (p *orderlist) FindByID(orderid string) Order {
	return p.ordermap[orderid]
}

func (p *orderlist) FindByPrice(price float64) []Order {
	var orders []Order

	for iter := p.orderedlist.Iterator(); iter.Next(); {
		if floatEquals(iter.Value().(Order).Price(), price) {
			order := iter.Value().(Order)
			orders = append(orders, order)
			fmt.Println("order", order)
		}
	}
	return orders
}

const epsilon float64 = 0.00000001

func floatEquals(a, b float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
