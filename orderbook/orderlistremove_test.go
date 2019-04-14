package orderbook

import (
	"testing"
)

func Test_RemoveOrder(m *testing.T) {
	pt := threeOrders(m, TopIsLow)

	pt.Remove("orderid1")
	if pt.Size() != 2 {
		m.Error("Size not 2 was", pt.Size())
	}

	orders := pt.GetAll()
	assertequal(m, orders[0].orderid, "orderid3", "RemoveOrder")
	assertequal(m, orders[1].orderid, "orderid2", "RemoveOrder")
	assertequal(m, orders[0].price, 1.0, "RemoveOrder")
	assertequal(m, orders[1].price, 1.1, "RemoveOrder")

	if pt.Remove("orderid1") != false {
		m.Errorf("Second remove should be false")
	}
	if pt.Size() != 2 {
		m.Error("Size not 2 was", pt.Size())
	}

	if pt.Remove("orderid2") != true {
		m.Errorf("orderid2 remove should be true")
	}
	if pt.Remove("orderid3") != true {
		m.Errorf("orderid3 remove should be true")
	}
	if pt.Size() != 0 {
		m.Error("Size not 0 was", pt.Size())
	}
	if pt.Remove("orderid1") != false {
		m.Errorf("third remove should be false")
	}

}
