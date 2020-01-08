package orderstate

import (
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/emirpasic/gods/utils"
	"github.com/google/uuid"
	"orderbook/etype"
)

type OrderOfList int

const (
	LowToHigh OrderOfList = -1 // Sell orders increase in price
	HighToLow             = 1  // Buy orders decrease in price
)

const (
	duplicateOrder = etype.Error("duplicate order")
)

type Order interface {
	OrderID() string
	Price() float64
	Timestamp() uuid.UUID
}

type Orderlist struct {
	orderedlist *treeset.Set
	ordermap    map[string]*OrderState
	comparator  utils.Comparator
}

func NewOrderList(comparator utils.Comparator) *Orderlist {
	p := Orderlist{}
	p.orderedlist = treeset.NewWith(comparator)
	p.comparator = comparator
	p.ordermap = make(map[string]*OrderState)
	return &p
}

func (p *Orderlist) Add(order *OrderState) error {
	if ord := p.ordermap[order.OrderID()]; ord != nil {
		return duplicateOrder
	}
	p.orderedlist.Add(order)
	p.ordermap[order.OrderID()] = order
	return nil
}

func (p *Orderlist) Size() int {
	return p.orderedlist.Size()
}

func (p *Orderlist) RemoveByID(orderid string) bool {
	if ord := p.FindByID(orderid); ord != nil {
		p.orderedlist.Remove(ord)
		delete(p.ordermap, orderid)
		return true
	}
	return false
}

func (p *Orderlist) Top() *OrderState {
	if p.orderedlist.Size() == 0 {
		return nil
	}
	var iter = p.orderedlist.Iterator()
	iter.Next()
	return iter.Value().(*OrderState)
}

func (p *Orderlist) Orders() []*OrderState {
	var orders []*OrderState
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(*OrderState)
		orders = append(orders, order)
	}
	return orders
}

func (p *Orderlist) Orders2(ord []*OrderState) []*OrderState {
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(*OrderState)
		ord = append(ord, order)
	}
	return ord
}

func (p *Orderlist) CopyList() *Orderlist {
	list := NewOrderList(p.comparator)
	for iter := p.Iterator(); iter.Next() == true; {
		err := list.Add(iter.Value().(*OrderState))
		if err != nil {
			panic("it should not happen")
		}
	}

	return list
}

func (p *Orderlist) FindByID(orderid string) *OrderState {
	return p.ordermap[orderid]
}

func (p *Orderlist) FindFirst(predicate func(order interface{}) bool) *OrderState {
	for iter := p.orderedlist.Iterator(); iter.Next(); {
		if predicate(iter.Value()) {
			return iter.Value().(*OrderState)
		}
	}
	return nil
}

func (p *Orderlist) FindAll(predicate func(order interface{}) bool) []*OrderState {
	var orders []*OrderState
	for iter := p.orderedlist.Iterator(); iter.Next(); {
		if predicate(iter.Value()) {
			order := iter.Value().(*OrderState)
			orders = append(orders, order)
		}
	}
	return orders
}

func (p *Orderlist) Iterator() treeset.Iterator {
	return p.orderedlist.Iterator()
}
