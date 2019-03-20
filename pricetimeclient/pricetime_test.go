package pricetimeclient

import (
	"orderbook/pricetime"
	"testing"
	"time"
)

func Test_AddThreeOrders(m *testing.T) {
	pt := pricetime.NewPriceTime()
	var err error
	err = pt.Add(pricetime.NewPriceTimeItem("orderid1", 1.2, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(pricetime.NewPriceTimeItem("orderid2", 1.1, time.Now(), "data"))
	printerror(err, m)
	err = pt.Add(pricetime.NewPriceTimeItem("orderid3", 1.1, time.Now(), "data"))
	printerror(err, m)

	if pt.Size() != 3 {
		m.Error("Size not 3 was", pt.Size())
	}

	if pti := pt.Top(); pti.Orderid() != "orderid2" {
		m.Error("Price Error", pti.Price(), pti.Orderid())
	}
}

func printerror(err error, m *testing.T) {
	if err != nil {
		m.Errorf("Add order failed %v", err)
	}
}
