package orderbook

import (
	"github.com/google/uuid"
	"testing"
)

func Test_RejectDuplicateOrderID(m *testing.T) {
	pt := threeOrders(m, TopIsLow)
	var err error
	err = pt.Add(NewOrder("orderid2", 1.1, newID(uuid.NewUUID()), "data"))
	if err != DuplicateOrder {
		m.Errorf("err should be %v was %v", DuplicateOrder, err)
	}
	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
	}

	if pti := pt.Top(); pti.Price() != 1.1 && pti.Orderid() == "orderid2" {
		m.Error("Price Error", pti.Price())
	}

	orders := pt.GetAll()

	assertequal(m, orders[0].orderid, "orderid3", "RejectDuplicateOrder")
	assertequal(m, orders[1].orderid, "orderid2", "RejectDuplicateOrder")
	assertequal(m, orders[2].orderid, "orderid1", "RejectDuplicateOrder")
}
