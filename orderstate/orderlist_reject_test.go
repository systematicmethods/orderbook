package orderstate

import (
	"github.com/google/uuid"
	"orderbook/assert"
	"testing"
)

func Test_RejectDuplicateOrderID(t *testing.T) {
	pt := threeOrders(t, LowToHigh)
	var err error
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuid.NewUUID())))
	if err != duplicateOrder {
		t.Errorf("err should be %v was %v", duplicateOrder, err)
	}
	if pt.Size() != 3 {
		t.Error("Size not 3 was", pt.Size())
	}

	if pti := pt.Top(); pti.Price() != 1.1 && pti.OrderID() == "orderid2" {
		t.Error("Price Error", pti.Price())
	}

	orders := pt.Orders()

	assert.AssertEqualT(t, orders[0].OrderID(), "orderid3", "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[1].OrderID(), "orderid2", "RejectDuplicateOrder")
	assert.AssertEqualT(t, orders[2].OrderID(), "orderid1", "RejectDuplicateOrder")
}
