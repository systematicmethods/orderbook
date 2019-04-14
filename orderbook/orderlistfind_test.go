package orderbook

import (
	"testing"
)

func Test_FindOrderByID(m *testing.T) {
	pt := threeOrders(m, TopIsLow)

	order := pt.FindByID("orderid1")

	if assertNotNil(m, order, "FindOrderByID") {
		assertequal(m, order.Orderid(), "orderid1", "FindOrderByID")
		assertequal(m, order.Price(), 1.2, "FindOrderByID")
	}
}

func Test_DidNotFindOrderByID(m *testing.T) {
	pt := threeOrders(m, TopIsLow)

	order := pt.FindByID("orderid11")

	assertNil(m, order, "FindOrderByID")

}

func Test_FindOrderByPrice(m *testing.T) {
	pt := threeOrders(m, TopIsLow)

	order := pt.FindByPrice(1.2)
	assertequal(m, len(order), 1, "num orders should be 1")

	if assertNotNil(m, order[0], "FindOrderByPrice") {
		assertequal(m, order[0].Orderid(), "orderid1", "FindOrderByPrice")
		assertequal(m, order[0].Price(), 1.2, "FindOrderByPrice")
	}
}

func Test_FindOrderByPriceWithTwoPrices(m *testing.T) {
	pt := threeOrdersTwoAtSamePrice(m, TopIsLow)

	order := pt.FindByPrice(1.1)
	assertequal(m, len(order), 2, "num orders should be 1")

	if assertNotNil(m, order[0], "FindOrderByPrice two prices") {
		assertequal(m, order[0].Orderid(), "orderid2", "FindOrderByPrice two prices")
		assertequal(m, order[0].Price(), 1.1, "FindOrderByPrice two prices")
	}

	if assertNotNil(m, order[1], "FindOrderByPrice two prices") {
		assertequal(m, order[1].Orderid(), "orderid3", "FindOrderByPrice two prices")
		assertequal(m, order[1].Price(), 1.1, "FindOrderByPrice two prices")
	}
}
