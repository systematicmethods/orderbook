package orderbook

import (
	"testing"
	"time"
)

func Test_AddThreeOrders(m *testing.T) {
	pt := NewPriceTime()
	var err error

	err = pt.Add(NewOrder("orderid1", 1.2, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid2", 1.1, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid3", 1.1, time.Now(), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
	}

	if pti := pt.Top(); pti.orderid != "orderid2" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	//UnixTime

	orders := pt.AllItems()

	orders[0].orderid = "orderid1"
	orders[1].orderid = "orderid2"
	orders[2].orderid = "orderid3"
}

func Test_RejectDuplicateOrderID(m *testing.T) {
	pt := NewPriceTime()
	var err error
	err = pt.Add(NewOrder("orderid1", 1.2, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid2", 1.1, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid2", 1.1, time.Now(), "data"))
	if err != DuplicateOrder {
		m.Errorf("err should be %v was %v", DuplicateOrder, err)
	}
	if pt.Size() != 2 {
		m.Error("Size not 2 was", pt.Size())
	}

	if pti := pt.Top(); pti.Price() != 1.1 && pti.Orderid() == "orderid2" {
		m.Error("Price Erro", pti.Price())
	}

	orders := pt.AllItems()

	orders[0].orderid = "orderid1"
	orders[1].orderid = "orderid2"
}

func Test_ReomveOrder(m *testing.T) {
	pt := NewPriceTime()
	var err error
	err = pt.Add(NewOrder("orderid1", 1.2, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid2", 1.1, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid3", 1.1, time.Now(), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
	}

	pt.Remove("orderid1")
	orders := pt.AllItems()

	orders[0].orderid = "orderid2"
	orders[1].orderid = "orderid3"
}

func printerror(err error, m *testing.T) {
	if err != nil {
		m.Errorf("Add order failed %v", err)
	}
}
