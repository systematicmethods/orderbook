package orderbook

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
)

var DuplicateOrder = errors.New("Duplicate Order")

type PriceTime interface {
	Add(order *order) *priceTime
	Top() *order
	Remove(orderid string)
	AllItems() []*order
}

type priceTime struct {
	orderedlist *treeset.Set
	itemset     map[string]bool
}

func NewPriceTime() *priceTime {
	p := priceTime{}
	p.orderedlist = treeset.NewWith(priceComparator)
	p.itemset = make(map[string]bool)
	return &p
}

func (p *priceTime) Add(order *order) error {
	if p.itemset[order.Orderid()] {
		return DuplicateOrder
	}
	p.orderedlist.Add(order)
	p.itemset[order.Orderid()] = true
	return nil
}

func (p *priceTime) Size() int {
	return p.orderedlist.Size()
}

func (p *priceTime) Remove(orderid string) {

}

func (p *priceTime) Top() *order {
	var iter = p.orderedlist.Iterator()
	iter.Next()
	return iter.Value().(*order)
}

func (p *priceTime) AllItems() []*order {
	var orders []*order
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(*order)
		fmt.Println("order", order)
		orders = append(orders, order)
	}
	return orders
}
