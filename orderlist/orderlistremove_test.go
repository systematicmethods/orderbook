package orderlist

import (
	"orderbook/assert"
	"testing"
)

func Test_RemoveOrder(m *testing.T) {
	pt := threeOrders(m, LowToHigh)

	pt.RemoveByID("orderid1")
	if pt.Size() != 2 {
		m.Error("Size not 2 was", pt.Size())
	}

	orders := pt.GetAll()
	assert.AssertEqual(m, orders[0].orderid, "orderid3", "RemoveOrder")
	assert.AssertEqual(m, orders[1].orderid, "orderid2", "RemoveOrder")
	assert.AssertEqual(m, orders[0].price, 1.0, "RemoveOrder")
	assert.AssertEqual(m, orders[1].price, 1.1, "RemoveOrder")

	if pt.RemoveByID("orderid1") != false {
		m.Errorf("Second remove should be false")
	}
	if pt.Size() != 2 {
		m.Error("Size not 2 was", pt.Size())
	}

	if pt.RemoveByID("orderid2") != true {
		m.Errorf("orderid2 remove should be true")
	}
	if pt.RemoveByID("orderid3") != true {
		m.Errorf("orderid3 remove should be true")
	}
	if pt.Size() != 0 {
		m.Error("Size not 0 was", pt.Size())
	}
	if pt.RemoveByID("orderid1") != false {
		m.Errorf("third remove should be false")
	}

}
