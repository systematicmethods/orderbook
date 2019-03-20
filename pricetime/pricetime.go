package pricetime

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/sets/treeset"
)

var DuplicateOrder = errors.New("Duplicate Order")

type PriceTime interface {
	Add(order *priceTimeItem) *priceTime
	Top() *priceTimeItem
	Remove(orderid string)
	AllItems() []*priceTimeItem
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

func (p *priceTime) Add(order *priceTimeItem) error {
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

func (p *priceTime) Top() *priceTimeItem {
	var iter = p.orderedlist.Iterator()
	iter.Next()
	return iter.Value().(*priceTimeItem)
}

func (p *priceTime) AllItems() []*priceTimeItem {
	var orders []*priceTimeItem
	for iter := p.orderedlist.Iterator(); iter.Next() == true; {
		order := iter.Value().(*priceTimeItem)
		fmt.Println("order", order)
		orders = append(orders, order)
	}
	return orders
}
