package orderstate

import (
	"orderbook/assert"
	"orderbook/test"
	"orderbook/uuidext"
	"testing"
	"time"
)

func Test_AddThreeOrdersSellSide(t *testing.T) {
	pt := threeOrdersTwoAtSamePrice(t, sellPriceComparator)

	if pti := pt.Top(); pti.OrderID() != "orderid2" {
		t.Error("Price Error", pti)
	}

	orders := pt.Orders()

	assert.AssertEqualT(t, orders[0].OrderID(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(t, orders[1].OrderID(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(t, orders[2].OrderID(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(t, orders[0].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[2].Price(), 1.2, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySide(t *testing.T) {
	pt := threeOrdersTwoAtSamePrice(t, buyPriceComparator)

	if pti := pt.Top(); pti.OrderID() != "orderid1" {
		t.Error("Price Error", pti.Price(), pti.OrderID())
	}

	orders := pt.Orders()

	assert.AssertEqualT(t, orders[0].OrderID(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(t, orders[1].OrderID(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(t, orders[2].OrderID(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(t, orders[0].Price(), 1.2, "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[2].Price(), 1.1, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySideGeneratedTime(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 100, loc)
	pt := NewOrderList(buyPriceComparator)

	var err error
	err = pt.Add(newOrderForTesting("clordid1", "orderid1", 1.2, newID(uuidext.NewUUIDFromTimeSeq(dt, 1))))
	test.PrintError(err, t)
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 2))))
	test.PrintError(err, t)
	err = pt.Add(newOrderForTesting("clordid3", "orderid3", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 3))))
	test.PrintError(err, t)

	if pt.Size() != 3 {
		t.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(t, pti.Timestamp(), pti.OrderID())
		}
	}

	if pti := pt.Top(); pti.OrderID() != "orderid1" {
		t.Error("Price Error", pti.Price(), pti.OrderID())
	}

	orders := pt.Orders()

	assert.AssertEqualT(t, orders[0].OrderID(), "orderid1", "AddThreeOrders")
	assert.AssertEqualT(t, orders[1].OrderID(), "orderid2", "AddThreeOrders")
	assert.AssertEqualT(t, orders[2].OrderID(), "orderid3", "AddThreeOrders")
	assert.AssertEqualT(t, orders[0].Price(), 1.2, "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[1].Price(), 1.1, "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[2].Price(), 1.1, "RejectDuplicateOrder")
}
