package orderbook

import (
	"github.com/google/uuid"
	"testing"
)

func threeOrders(m *testing.T, orderby OrderOfList) *orderlist {
	pt := NewOrderListStruct(orderby)
	var err error
	err = pt.Add(newOrderForTesting("clordid1", "orderid1", 1.2, newID(uuid.NewUUID()), "data"))
	printerror(err, m)
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuid.NewUUID()), "data"))
	printerror(err, m)
	err = pt.Add(newOrderForTesting("clordid3", "orderid3", 1.0, newID(uuid.NewUUID()), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(m, pti.Timestamp(), pti.OrderID())
		}
	}

	return pt
}

func threeOrdersTwoAtSamePrice(m *testing.T, orderby OrderOfList) *orderlist {
	pt := NewOrderListStruct(orderby)
	var err error
	err = pt.Add(newOrderForTesting("clordid1", "orderid1", 1.2, newID(uuid.NewUUID()), "data"))
	printerror(err, m)
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuid.NewUUID()), "data"))
	printerror(err, m)
	err = pt.Add(newOrderForTesting("clordid3", "orderid3", 1.1, newID(uuid.NewUUID()), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(m, pti.Timestamp(), pti.OrderID())
		}
	}

	return pt
}
