package orderbook

import (
	"orderbook/assert"
	"orderbook/uuidext"
	"testing"
	"time"
)

func Test_AddThreeOrdersSellSide(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, LowToHigh)

	if pti := pt.Top(); pti.OrderID() != "orderid2" {
		m.Error("Price Error", pti.Price(), pti.OrderID())
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].OrderID(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(m, orders[1].OrderID(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(m, orders[2].OrderID(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(m, orders[0].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].Price(), 1.2, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySide(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, HighToLow)

	if pti := pt.Top(); pti.OrderID() != "orderid1" {
		m.Error("Price Error", pti.Price(), pti.OrderID())
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].OrderID(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(m, orders[1].OrderID(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(m, orders[2].OrderID(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(m, orders[0].Price(), 1.2, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].Price(), 1.1, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySideGeneratedTime(m *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 100, loc)
	pt := NewOrderListStruct(HighToLow)
	var err error

	err = pt.Add(newOrderForTesting("clordid1", "orderid1", 1.2, newID(uuidext.NewUUIDFromTimeSeq(dt, 1)), "data"))
	printerror(err, m)
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 2)), "data"))
	printerror(err, m)
	err = pt.Add(newOrderForTesting("clordid3", "orderid3", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 3)), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(m, pti.Timestamp(), pti.OrderID())
		}
	}

	if pti := pt.Top(); pti.OrderID() != "orderid1" {
		m.Error("Price Error", pti.Price(), pti.OrderID())
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].OrderID(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(m, orders[1].OrderID(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(m, orders[2].OrderID(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(m, orders[0].Price(), 1.2, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].Price(), 1.1, "RejectDuplicateOrder")
}
