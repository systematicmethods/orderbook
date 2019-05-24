package orderbook

import (
	"github.com/google/uuid"
	"orderbook/assert"
	"testing"
)

func Test_RejectDuplicateOrderID(m *testing.T) {
	pt := threeOrders(m, LowToHigh)
	var err error
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuid.NewUUID()), "data"))
	if err != DuplicateOrder {
		m.Errorf("err should be %v was %v", DuplicateOrder, err)
	}
	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
	}

	if pti := pt.Top(); pti.Price() != 1.1 && pti.OrderID() == "orderid2" {
		m.Error("Price Error", pti.Price())
	}

	orders := pt.Orders()

	assert.AssertEqualT(m, orders[0].OrderID(), "orderid3", "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[1].OrderID(), "orderid2", "RejectDuplicateOrder")
	assert.AssertEqualT(m, orders[2].OrderID(), "orderid1", "RejectDuplicateOrder")
}
