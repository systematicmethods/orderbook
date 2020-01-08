package orderstate

import (
	"fmt"
)

func printOrders(lst *Orderlist) {
	fmt.Printf("clientid|clordid|side|price|qty|ordstatus \n")
	for iter := lst.Iterator(); iter.Next() == true; {
		order := iter.Value().(OrderState)
		fmt.Printf("%s|%s|%v|%v|%v \n", order.ClientID(), order.ClOrdID(), order.Side(), order.OrderQty(), order.Price())
	}
}
