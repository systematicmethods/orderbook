package orderstate

import (
	"fmt"
	"orderbook/assert"
	"testing"
)

func Test_FindOrderByID(t *testing.T) {
	pt := threeOrders(t, LowToHigh)

	order := pt.FindByID("orderid1")

	if assert.AssertNotNilT(t, order, "FindOrderByID") {
		assert.AssertEqualT(t, order.OrderID(), "orderid1", "FindOrderByID")
		assert.AssertEqualT(t, order.Price(), 1.2, "FindOrderByID")
	}
}

func Test_DidNotFindOrderByID(t *testing.T) {
	pt := threeOrders(t, LowToHigh)
	order := pt.FindByID("orderid11")
	fmt.Printf("order %v", order)
	assert.AssertNilT(t, order, "FindOrderByID")
}

func Test_FindOrderByPrice(t *testing.T) {
	pt := threeOrders(t, LowToHigh)

	order := pt.FindAll(func(ord interface{}) bool {
		return ord.(Order).Price() == 1.2
	})

	assert.AssertEqualT(t, len(order), 1, "num orders should be 1")

	if assert.AssertNotNilT(t, order[0], "FindOrderByPrice") {
		assert.AssertEqualT(t, order[0].OrderID(), "orderid1", "FindOrderByPrice")
		assert.AssertEqualT(t, order[0].Price(), 1.2, "FindOrderByPrice")
	}
}

func Test_FindOrderByPriceWithTwoPrices(t *testing.T) {
	pt := threeOrdersTwoAtSamePrice(t, sellPriceComparator)

	order := pt.FindAll(func(ord interface{}) bool {
		return ord.(Order).Price() == 1.1
	})

	assert.AssertEqualT(t, len(order), 2, "num orders should be 1")

	if len(order) > 0 && assert.AssertNotNilT(t, order[0], "FindOrderByPrice two prices") {
		assert.AssertEqualT(t, order[0].OrderID(), "orderid2", "FindOrderByPrice two prices")
		assert.AssertEqualT(t, order[0].Price(), 1.1, "FindOrderByPrice two prices")
	}

	if len(order) > 1 && assert.AssertNotNilT(t, order[1], "FindOrderByPrice two prices") {
		assert.AssertEqualT(t, order[1].OrderID(), "orderid3", "FindOrderByPrice two prices")
		assert.AssertEqualT(t, order[1].Price(), 1.1, "FindOrderByPrice two prices")
	}
}
