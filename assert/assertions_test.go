package assert

import "testing"

func Test_RejectDuplicateOrderID(t *testing.T) {

	AssertEqualT(t, "orderid3", "orderid3", "equal")
	AssertEqual("orderid3", "orderid3", "equal")
	AssertNotEqualT(t, "orderid3", "orderid5", "not equal")
	AssertNotNilT(t, "orderid5", "not nil")
	AssertNilT(t, nil, "not nil")

}
