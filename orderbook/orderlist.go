package orderbook

import (
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
)

type OrderOfList int

const (
	TopIsLow  OrderOfList = -1 // Sell orders increase in price
	TopIsHigh             = 1  // Buy orders decrease in price
)

type OrderList interface {
	Add(order *order) *orderlist
	Top() *order
	RemoveByID(orderid string)
	FindByID(orderid string) *order
	FindByPrice(price float64) []*order
	GetAll() []*order
}

type orderlist struct {
	orderedlist *treeset.Set
	ordermap    map[string]*order
}

func NewOrderList(side OrderOfList) *orderlist {
	p := orderlist{}
	if side == TopIsLow {
		p.orderedlist = treeset.NewWith(sellPriceComparator)
	} else if side == TopIsHigh {
		p.orderedlist = treeset.NewWith(buyPriceComparator)
	}
	p.ordermap = make(map[string]*order)
	return &p
}

func (p *orderlist) Add(order *order) error {
	if ord := p.ordermap[order.Orderid()]; ord != nil {
		return DuplicateOrder
	}
	p.orderedlist.Add(order)
	p.ordermap[order.Orderid()] = order
	return nil
}

func (p *orderlist) Size() int {
	return p.orderedlist.Size()
}

func (p *orderlist) Remove(orderid string) bool {
	if ord := p.FindByID(orderid); ord != nil {
		p.orderedlist.Remove(ord)
		delete(p.ordermap, orderid)
		return true
	}
	return false
}

func (p *orderlist) Top() *order {
	var iter = p.orderedlist.Iterator()
	iter.Next()
	return iter.Value().(*order)
}

func (p *orderlist) GetAll() []*order {
	var orders []*order
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(*order)
		orders = append(orders, order)
	}
	return orders
}

func (p *orderlist) FindByID(orderid string) *order {
	return p.ordermap[orderid]
}

func (p *orderlist) FindByPrice(price float64) []*order {
	var orders []*order

	for iter := p.orderedlist.Iterator(); iter.Next(); {
		if floatEquals(iter.Value().(*order).price, price) {
			order := iter.Value().(*order)
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
