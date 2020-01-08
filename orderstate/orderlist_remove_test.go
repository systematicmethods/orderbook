package orderstate

import (
	"orderbook/assert"
	"testing"
)

func Test_RemoveOrder(t *testing.T) {
	pt := threeOrders(t, LowToHigh)

	pt.RemoveByID("orderid1")
	if pt.Size() != 2 {
		t.Error("Size not 2 was", pt.Size())
	}

	orders := pt.Orders()
	assert.AssertEqualT(t, orders[0].OrderID(), "orderid3", "RemoveOrder")
	assert.AssertEqualT(t, orders[1].OrderID(), "orderid2", "RemoveOrder")
	assert.AssertEqualT(t, orders[0].Price(), 1.0, "RemoveOrder")
	assert.AssertEqualT(t, orders[1].Price(), 1.1, "RemoveOrder")

	if pt.RemoveByID("orderid1") != false {
		t.Errorf("Second remove should be false")
	}
	if pt.Size() != 2 {
		t.Error("Size not 2 was", pt.Size())
	}

	if pt.RemoveByID("orderid2") != true {
		t.Errorf("orderid2 remove should be true")
	}
	if pt.RemoveByID("orderid3") != true {
		t.Errorf("orderid3 remove should be true")
	}
	if pt.Size() != 0 {
		t.Error("Size not 0 was", pt.Size())
	}
	if pt.RemoveByID("orderid1") != false {
		t.Errorf("third remove should be false")
	}

}
