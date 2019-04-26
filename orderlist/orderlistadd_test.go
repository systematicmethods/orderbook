package orderlist

import (
	"orderbook/assert"
	"orderbook/uuidext"
	"testing"
	"time"
)

func Test_AddThreeOrdersSellSide(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, LowToHigh)

	if pti := pt.Top(); pti.orderid != "orderid2" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].Orderid(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(m, orders[1].Orderid(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(m, orders[2].Orderid(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(m, orders[0].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].Price(), 1.2, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySide(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, HighToLow)

	if pti := pt.Top(); pti.orderid != "orderid1" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].Orderid(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(m, orders[1].Orderid(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(m, orders[2].Orderid(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(m, orders[0].Price(), 1.2, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].Price(), 1.1, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySideGeneratedTime(m *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 100, loc)
	pt := NewOrderListStruct(HighToLow)
	var err error

	err = pt.Add(NewOrder("orderid1", 1.2, newID(uuidext.NewUUIDFromTimeSeq(dt, 1)), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid2", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 2)), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid3", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 3)), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(m, pti.UUID(), pti.Orderid())
		}
	}

	if pti := pt.Top(); pti.orderid != "orderid1" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].Orderid(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(m, orders[1].Orderid(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(m, orders[2].Orderid(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(m, orders[0].Price(), 1.2, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].Price(), 1.1, "RejectDuplicateOrder")
}
