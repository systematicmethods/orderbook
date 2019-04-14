package orderbook

import (
	"orderbook/uuidext"
	"testing"
	"time"
)

func Test_AddThreeOrdersSellSide(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, TopIsLow)

	if pti := pt.Top(); pti.orderid != "orderid2" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	orders := pt.GetAll()

	assertequal(m, orders[0].orderid, "orderid2", "AddThreeOrders")
	assertequal(m, orders[1].orderid, "orderid3", "AddThreeOrders")
	assertequal(m, orders[2].orderid, "orderid1", "AddThreeOrders")
	assertequal(m, orders[0].price, 1.1, "RejectDuplicateOrder")
	assertequal(m, orders[1].price, 1.1, "RejectDuplicateOrder")
	assertequal(m, orders[2].price, 1.2, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySide(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, TopIsHigh)

	if pti := pt.Top(); pti.orderid != "orderid1" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	orders := pt.GetAll()

	assertequal(m, orders[0].orderid, "orderid1", "AddThreeOrders")
	assertequal(m, orders[1].orderid, "orderid2", "AddThreeOrders")
	assertequal(m, orders[2].orderid, "orderid3", "AddThreeOrders")
	assertequal(m, orders[0].price, 1.2, "RejectDuplicateOrder")
	assertequal(m, orders[1].price, 1.1, "RejectDuplicateOrder")
	assertequal(m, orders[2].price, 1.1, "RejectDuplicateOrder")
}

func Test_AddThreeOrdersBuySideGeneratedTime(m *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	dt := time.Date(2019, 10, 11, 11, 11, 1, 100, loc)
	pt := NewOrderList(TopIsHigh)
	var err error

	err = pt.Add(NewOrder("orderid1", 1.2, newID(uuidext.NewUUIDFromTimeSeq(dt, 1)), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid2", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 2)), "data"))
	printerror(err, m)
	err = pt.Add(NewOrder("orderid3", 1.1, newID(uuidext.NewUUIDFromTimeSeq(dt, 3)), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.GetAll() {
			dumptime(m, pti.timeuuid, pti.orderid)
		}
	}

	if pti := pt.Top(); pti.orderid != "orderid1" {
		m.Error("Price Error", pti.price, pti.orderid)
	}

	orders := pt.GetAll()

	assertequal(m, orders[0].orderid, "orderid1", "AddThreeOrders")
	assertequal(m, orders[1].orderid, "orderid2", "AddThreeOrders")
	assertequal(m, orders[2].orderid, "orderid3", "AddThreeOrders")
	assertequal(m, orders[0].price, 1.2, "RejectDuplicateOrder")
	assertequal(m, orders[1].price, 1.1, "RejectDuplicateOrder")
	assertequal(m, orders[2].price, 1.1, "RejectDuplicateOrder")
}
