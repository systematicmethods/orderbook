package orderstate

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"orderbook/test"
	"orderbook/uuidext"
	"testing"
)

func sellPriceComparator(a, b interface{}) int {
	apti := a.(*OrderState)
	bpti := b.(*OrderState)
	switch {
	case apti.price > bpti.price:
		return 1
	case apti.price < bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}

func buyPriceComparator(a, b interface{}) int {
	apti := a.(*OrderState)
	bpti := b.(*OrderState)
	switch {
	case apti.price < bpti.price:
		return 1
	case apti.price > bpti.price:
		return -1
	default:
		return uuidext.UUIDComparator(apti.timestamp, bpti.timestamp)
	}
}

func threeOrders(t *testing.T, orderby OrderOfList) *Orderlist {
	var pt *Orderlist
	if orderby == LowToHigh {
		pt = NewOrderList(sellPriceComparator)
	} else {
		pt = NewOrderList(buyPriceComparator)
	}
	var err error
	err = pt.Add(newOrderForTesting("clordid1", "orderid1", 1.2, newID(uuid.NewUUID())))
	test.PrintError(err, t)
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuid.NewUUID())))
	test.PrintError(err, t)
	err = pt.Add(newOrderForTesting("clordid3", "orderid3", 1.0, newID(uuid.NewUUID())))
	test.PrintError(err, t)

	if pt.Size() != 3 {
		t.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(t, pti.Timestamp(), pti.OrderID())
		}
	}

	return pt
}

func threeOrdersTwoAtSamePrice(t *testing.T, f func(a, b interface{}) int) *Orderlist {
	pt := NewOrderList(f)
	var err error
	err = pt.Add(newOrderForTesting("clordid1", "orderid1", 1.2, newID(uuid.NewUUID())))
	test.PrintError(err, t)
	err = pt.Add(newOrderForTesting("clordid2", "orderid2", 1.1, newID(uuid.NewUUID())))
	test.PrintError(err, t)
	err = pt.Add(newOrderForTesting("clordid3", "orderid3", 1.1, newID(uuid.NewUUID())))
	test.PrintError(err, t)

	if pt.Size() != 3 {
		t.Error("Size not 3 was", pt.Size())
		for _, pti := range pt.Orders() {
			dumptime(t, pti.Timestamp(), pti.OrderID())
		}
	}

	return pt
}

func newOrderForTesting(clOrdID string, orderID string, price float64, timestamp uuid.UUID) *OrderState {
	return &OrderState{clOrdID: clOrdID, orderID: orderID, price: price, timestamp: timestamp}
}

func dumptime(t *testing.T, id uuid.UUID, msg string) {
	t.Errorf("Time %s %d, %d, %v %s", msg, id.Time(), id.ClockSequence(), id.Version(), hex.Dump(id[:]))
	dumpbytes(id[:])
}

func dumpbytes(b []byte) {
	for _, n := range b[:] {
		fmt.Printf(" %08b", n) // prints 00000000 11111101
	}
	fmt.Println()
}

func newID(uuid uuid.UUID, _ error) uuid.UUID {
	return uuid
}
