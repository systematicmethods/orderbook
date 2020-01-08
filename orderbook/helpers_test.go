package orderbook

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"orderbook/fixmodel"
	"orderbook/orderstate"
	"testing"
)

const inst = "ABV"

func dumptime(m *testing.T, id uuid.UUID, msg string) {
	m.Errorf("Time %s %d, %d, %v %s", msg, id.Time(), id.ClockSequence(), id.Version(), hex.Dump(id[:]))
	dumpbytes(id[:])
}

func dumpbytes(b []byte) {
	for _, n := range b[:] {
		fmt.Printf(" %08b", n) // prints 00000000 11111101
	}
	fmt.Println()
}

func printExecsAndOrders(execs []fixmodel.ExecutionReport, bk *BuySellOrders, buyorders []orderstate.OrderState, sellorders []orderstate.OrderState) {
	for i, s := range execs {
		fmt.Printf("e%d %v\n", i, s)
	}
	fmt.Printf("id|clientid|clordid|side|lastprice|lastqty|status|price|qty|ordstatus\n")
	for i, ex := range execs {
		var order *orderstate.OrderState
		if ex.Side() == fixmodel.SideBuy {
			for _, anorder := range buyorders {
				if anorder.ClOrdID() == ex.ClOrdID() {
					order = &anorder
				}
			}
		} else {
			for _, anorder := range sellorders {
				if anorder.ClOrdID() == ex.ClOrdID() {
					order = &anorder
				}
			}
		}
		if order != nil {
			fmt.Printf("e%d|%s|%s|%s|%v|%v|%v|%v|%v|%v\n", i, ex.ClientID(), ex.ClOrdID(), ex.Side(), ex.LastPrice(), ex.LastQty(), ex.OrdStatus(), order.Price(), order.OrderQty(), order.OrdStatus())
		} else {
			fmt.Printf("e%d|%s|%s|%s|%v|%v|%v|\n", i, ex.ClientID(), ex.ClOrdID(), ex.Side(), ex.LastPrice(), ex.LastQty(), ex.OrdStatus())
		}
	}
}
