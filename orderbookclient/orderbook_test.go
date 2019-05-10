package orderbookclient

import (
	"fmt"
	"github.com/google/uuid"
	"orderbook/orderbook"
	"testing"
)

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}

func Test_AddThreeOrders(m *testing.T) {
	pt := orderbook.NewOrderList(orderbook.LowToHigh)
	var err error
	err = pt.Add(orderbook.NewOrder("orderid1", 1.2, newID(uuid.NewUUID()), "data"))
	printerror(err, m)
	err = pt.Add(orderbook.NewOrder("orderid2", 1.1, newID(uuid.NewUUID()), "data"))
	printerror(err, m)
	err = pt.Add(orderbook.NewOrder("orderid3", 1.1, newID(uuid.NewUUID()), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
	}

	if pti := pt.Top(); pti.Orderid() != "orderid2" {
		m.Error("Price Error", pti.Price(), pti.Orderid())
	}
}

func Test_AddLotsOfOrdersAtSamePriceLevel(m *testing.T) {
	exs := []struct {
		id  uuid.UUID
		oid string
	}{
		{newID(uuid.NewUUID()), "oid0"},
		{newID(uuid.NewUUID()), "oid1"},
		{newID(uuid.NewUUID()), "oid2"},
		{newID(uuid.NewUUID()), "oid3"},
		{newID(uuid.NewUUID()), "oid4"},
		{newID(uuid.NewUUID()), "oid5"},
		{newID(uuid.NewUUID()), "oid6"},
		{newID(uuid.NewUUID()), "oid7"},
		{newID(uuid.NewUUID()), "oid8"},
		{newID(uuid.NewUUID()), "oid9"},
		{newID(uuid.NewUUID()), "oid10"},
	}

	pt := orderbook.NewOrderList(orderbook.LowToHigh)
	var err error
	for _, auuid := range exs {
		err = pt.Add(orderbook.NewOrder(auuid.oid, 1.2, auuid.id, "data"))
		printerror(err, m)
	}

	if pt.Size() != len(exs) {
		m.Errorf("Size not %d was %d", len(exs), pt.Size())
	}

	orders := pt.Orders()
	for i := 0; i < 10; i++ {
		assertequal(m, orders[i].Orderid(), fmt.Sprintf("oid%d", i), "AddLotsOfOrdersAtSamePriceLevel")
	}
}

func assertequal(m *testing.T, a interface{}, b interface{}, msg string) {
	if a != b {
		m.Errorf("%s %s != %s", msg, a, b)
	}
}

func printerror(err error, m *testing.T) {
	if err != nil {
		m.Errorf("Add order failed %v", err)
	}
}
