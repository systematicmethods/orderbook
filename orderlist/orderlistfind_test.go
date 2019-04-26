package orderlist

import (
	"orderbook/assert"
	"testing"
)

func Test_FindOrderByID(m *testing.T) {
	pt := threeOrders(m, LowToHigh)

	order := pt.FindByID("orderid1")

	if assert.AssertNotNilT(m, order, "FindOrderByID") {
		assert.AssertEqualT(m, order.Orderid(), "orderid1", "FindOrderByID")
		assert.AssertEqualT(m, order.Price(), 1.2, "FindOrderByID")
	}
}

func Test_DidNotFindOrderByID(m *testing.T) {
	pt := threeOrders(m, LowToHigh)

	order := pt.FindByID("orderid11")

	assert.AssertNilT(m, order, "FindOrderByID")

}

func Test_FindOrderByPrice(m *testing.T) {
	pt := threeOrders(m, LowToHigh)

	order := pt.FindByPrice(1.2)
	assert.AssertEqualT(m, len(order), 1, "num orders should be 1")

	if assert.AssertNotNilT(m, order[0], "FindOrderByPrice") {
		assert.AssertEqualT(m, order[0].Orderid(), "orderid1", "FindOrderByPrice")
		assert.AssertEqualT(m, order[0].Price(), 1.2, "FindOrderByPrice")
	}
}

func Test_FindOrderByPriceWithTwoPrices(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, LowToHigh)

	order := pt.FindByPrice(1.1)
	assert.AssertEqualT(m, len(order), 2, "num orders should be 1")

	if assert.AssertNotNilT(m, order[0], "FindOrderByPrice two prices") {
		assert.AssertEqualT(m, order[0].Orderid(), "orderid2", "FindOrderByPrice two prices")
		assert.AssertEqualT(m, order[0].Price(), 1.1, "FindOrderByPrice two prices")
	}

	if assert.AssertNotNilT(m, order[1], "FindOrderByPrice two prices") {
		assert.AssertEqualT(m, order[1].Orderid(), "orderid3", "FindOrderByPrice two prices")
		assert.AssertEqualT(m, order[1].Price(), 1.1, "FindOrderByPrice two prices")
	}
}
