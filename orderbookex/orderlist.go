package orderbookex

import (
	"errors"
	"github.com/emirpasic/gods/sets/treeset"
)

type OrderOfList int

const (
	LowToHigh OrderOfList = -1 // Sell orders increase in price
	HighToLow             = 1  // Buy orders decrease in price
)

type OrderList interface {
	Add(order OrderState) error
	Top() (OrderState, error)
	RemoveByID(orderid string) bool
	FindByID(orderid string) OrderState
	FindByClOrdID(clOrdID string) OrderState
	FindByPrice(price float64) []OrderState
	Orders() []OrderState
	Size() int
	orderList() orderlist

	iterator() treeset.Iterator
}

type orderlist struct {
	orderedlist *treeset.Set
	ordermap    map[string]OrderState
}

func NewOrderListStruct(sort OrderOfList) *orderlist {
	p := orderlist{}
	if sort == LowToHigh {
		p.orderedlist = treeset.NewWith(sellPriceComparator)
	} else if sort == HighToLow {
		p.orderedlist = treeset.NewWith(buyPriceComparator)
	}
	p.ordermap = make(map[string]OrderState)
	return &p
}

func NewOrderList(sort OrderOfList) OrderList {
	return NewOrderListStruct(sort)
}

func (p *orderlist) Add(order OrderState) error {
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

func (p *orderlist) Top() (OrderState, error) {
	if p.orderedlist.Size() == 0 {
		var ret OrderState
		return ret, errors.New("No orders")
	}
	var iter = p.orderedlist.Iterator()
	iter.Next()
	return iter.Value().(OrderState), nil
}

func (p *orderlist) Orders() []OrderState {
	var orders []OrderState
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		orders = append(orders, order)
	}
	return orders
}

func (p *orderlist) orderList() orderlist {
	return *p
}

func (p *orderlist) FindByID(orderid string) OrderState {
	return p.ordermap[orderid]
}

func (p *orderlist) FindByClOrdID(clOrdID string) OrderState {
	for iter := p.orderedlist.Iterator(); iter.Next(); {
		if iter.Value().(OrderState).ClOrdID() == clOrdID {
			return iter.Value().(OrderState)
		}
	}
	return nil
}

func (p *orderlist) FindByPrice(price float64) []OrderState {
	var orders []OrderState

	for iter := p.orderedlist.Iterator(); iter.Next(); {
		if floatEquals(iter.Value().(OrderState).Price(), price) {
			order := iter.Value().(OrderState)
			orders = append(orders, order)
			//fmt.Println("order", order)
		}
	}
	return orders
}

func (p *orderlist) iterator() treeset.Iterator {
	return p.orderedlist.Iterator()
}

const epsilon float64 = 0.00000001

func floatEquals(a, b float64) bool {
	return (a-b) < epsilon && (b-a) < epsilon
}
